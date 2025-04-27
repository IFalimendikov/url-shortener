package transport

import (
    "bytes"
    "context"
    "encoding/json"
    "io"
    "net/http"
    "net/http/httptest"
    "os"
    "strings"
    "testing"
    "url-shortener/internal/config"
    "url-shortener/internal/handler"
    "url-shortener/internal/logger"
    "url-shortener/internal/models"
    "url-shortener/internal/services"
    "url-shortener/internal/storage"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestGetURL(t *testing.T) {
    t.Cleanup(func() {
        os.RemoveAll("test_storage")
    })

    cfg := config.Config{
        BaseURL:     "http://localhost:8080",
        StoragePath: "test_storage",
    }
    log := logger.NewLogger()
    storage, _ := storage.NewStorage(context.Background(), &cfg)
    defer storage.File.Close()

    service := services.NewURLService(context.Background(), log, storage)
    h := handler.NewHandler(service, log)
    tr := NewTransport(cfg, h, log)

    ts := httptest.NewServer(NewRouter(tr))
    defer ts.Close()

    client := ts.Client()
    client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse
    }

    tests := []struct {
        name       string
        shortURL   string
        longURL    string
        wantStatus int
        wantURL    string
    }{
        {
            name:       "Valid URL",
            shortURL:   "test1",
            longURL:    "https://example.com",
            wantStatus: http.StatusTemporaryRedirect,
            wantURL:    "https://example.com",
        },
        {
            name:       "Non-existent URL",
            shortURL:   "nonexistent",
            wantStatus: http.StatusBadRequest,
            wantURL:    "URL not found!",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.longURL != "" {
                userID := "test-user"
                shortURL, err := service.SaveURL(context.Background(), tt.longURL, userID)
                require.NoError(t, err)
                tt.shortURL = shortURL
            }

            resp, body := testRequest(t, ts, "GET", "/"+tt.shortURL, nil)
            defer resp.Body.Close()

            assert.Equal(t, tt.wantStatus, resp.StatusCode)
            if tt.wantStatus == http.StatusTemporaryRedirect {
                assert.Equal(t, tt.wantURL, resp.Header.Get("Location"))
            } else {
                assert.Equal(t, tt.wantURL, body)
            }
        })
    }
}

func TestPostURL(t *testing.T) {
    t.Cleanup(func() {
        os.RemoveAll("test_storage")
    })

    cfg := config.Config{
        BaseURL:     "http://localhost:8080",
        StoragePath: "test_storage",
    }
    log := logger.NewLogger()
    storage, _ := storage.NewStorage(context.Background(), &cfg)
    defer storage.File.Close()

    service := services.NewURLService(context.Background(), log, storage)
    h := handler.NewHandler(service, log)
    tr := NewTransport(cfg, h, log)

    ts := httptest.NewServer(NewRouter(tr))
    defer ts.Close()

    tests := []struct {
        name       string
        url        string
        wantStatus int
    }{
        {
            name:       "Valid URL",
            url:        "https://example.com",
            wantStatus: http.StatusCreated,
        },
        {
            name:       "Empty URL",
            url:        "",
            wantStatus: http.StatusBadRequest,
        },
        {
            name:       "Invalid URL",
            url:        "not-a-url",
            wantStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            resp, _ := testRequest(t, ts, "POST", "/", strings.NewReader(tt.url))
            defer resp.Body.Close()
            assert.Equal(t, tt.wantStatus, resp.StatusCode)
        })
    }
}

func TestShortenURL(t *testing.T) {
    t.Cleanup(func() {
        os.RemoveAll("test_storage")
    })

    cfg := config.Config{
        BaseURL:     "http://localhost:8080",
        StoragePath: "test_storage",
    }
    log := logger.NewLogger()
    storage, _ := storage.NewStorage(context.Background(), &cfg)
    defer storage.File.Close()

    service := services.NewURLService(context.Background(), log, storage)
    h := handler.NewHandler(service, log)
    tr := NewTransport(cfg, h, log)

    ts := httptest.NewServer(NewRouter(tr))
    defer ts.Close()

    tests := []struct {
        name       string
        request    models.ShortenURLRequest
        wantStatus int
    }{
        {
            name: "Valid Request",
            request: models.ShortenURLRequest{
                URL: "https://example.com",
            },
            wantStatus: http.StatusCreated,
        },
        {
            name: "Empty URL",
            request: models.ShortenURLRequest{
                URL: "",
            },
            wantStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            jsonBody, err := json.Marshal(tt.request)
            require.NoError(t, err)

            resp, _ := testRequest(t, ts, "POST", "/api/shorten", bytes.NewReader(jsonBody))
            defer resp.Body.Close()
            assert.Equal(t, tt.wantStatus, resp.StatusCode)
        })
    }
}

func TestPingDB(t *testing.T) {
    cfg := config.Config{
        BaseURL:     "http://localhost:8080",
        StoragePath: "test_storage",
    }
    log := logger.NewLogger()
    storage, _ := storage.NewStorage(context.Background(), &cfg)
    defer storage.File.Close()

    service := services.NewURLService(context.Background(), log, storage)
    h := handler.NewHandler(service, log)
    tr := NewTransport(cfg, h, log)

    ts := httptest.NewServer(NewRouter(tr))
    defer ts.Close()

    resp, body := testRequest(t, ts, "GET", "/ping", nil)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusOK, resp.StatusCode)
    assert.Equal(t, "Live", body)
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
    req, err := http.NewRequest(method, ts.URL+path, body)
    require.NoError(t, err)

    if method == "POST" && path == "/api/shorten" {
        req.Header.Set("Content-Type", "application/json")
    }

    resp, err := ts.Client().Do(req)
    require.NoError(t, err)

    respBody, err := io.ReadAll(resp.Body)
    require.NoError(t, err)

    return resp, string(respBody)
}
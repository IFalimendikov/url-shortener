package transport

import (
	"bytes"
	"io"
	"os"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"url-shortener/internal/config"
	"url-shortener/internal/logger"
	"url-shortener/internal/services"
	"url-shortener/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestPostURL(t *testing.T) {
	t.Cleanup(func() {
		os.RemoveAll("test_storage")
	})

	cfg := config.Config{
		BaseURL: "http://localhost:8080",
		StoragePath: "test_storage",
	}

	log := logger.NewLogger()
	storage, _ := storage.NewStorage(&cfg)
	defer storage.File.Close()

	s := services.NewURLService(log, storage)

	t.Cleanup(func() {
		os.RemoveAll("test_storage")
	})

	tr := NewTransport(cfg, s, log)

	ts := httptest.NewServer(NewRouter(cfg, tr))
	defer ts.Close()

	var testTable = []struct {
		url    string
		body   string
		want   string
		status int
	}{
		{"/", "https://practicum.yandex.ru/", "http://localhost:8080/5HZTcSScKLEKuTfRTUDifWUbO3kVswOjWFiZSB", http.StatusCreated},
		{"/", "https://practicum.yandex.at/", "http://localhost:8080/5HZTcSScKLEKuTfRTUDifWUbO3kVswOjWFdtYV", http.StatusCreated},
		{"/", "", "Empty body!", http.StatusBadRequest},
		{"/", "practicum.yandex.ru/", "Malformed URI!", http.StatusBadRequest},
	}

	for _, test := range testTable {
		resp, body := testRequest(t, ts, "POST", test.url, bytes.NewReader([]byte(test.body)))
		defer resp.Body.Close()
		assert.Equal(t, test.status, resp.StatusCode)
		assert.Equal(t, test.want, body)
	}
}

func TestGetURL(t *testing.T) {
	t.Cleanup(func() {
		os.RemoveAll("test_storage1")
	})

	cfg := config.Config{
		BaseURL: "http://localhost:8080",
		StoragePath: "test_storage1",
	}

	log := logger.NewLogger()
	storage, _ := storage.NewStorage(&cfg)
	defer storage.File.Close()
	s := services.NewURLService(log, storage)

	tr := NewTransport(cfg, s, log)

	ts := httptest.NewServer(NewRouter(cfg, tr))
	defer ts.Close()

	client := ts.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	var testTable = []struct {
		url    string
		body   string
		want   string
		status int
	}{
		{"/5HZTcSScKLEKuTfRTUDifWUbO3kVswOjWFiZSB", "", "https://practicum.yandex.ru/", http.StatusTemporaryRedirect},
		{"/5HZTcSScKLEKuTfRTUDifWUbO3kVswOjWFdtYV", "", "https://practicum.yandex.at/", http.StatusTemporaryRedirect},
		{"/3", "", "URL not found!", http.StatusBadRequest},
	}

	// First create the shortened URLs
	resp, _ := testRequest(t, ts, "POST", "/", bytes.NewReader([]byte("https://practicum.yandex.ru/")))
	defer resp.Body.Close()
	resp, _ = testRequest(t, ts, "POST", "/", bytes.NewReader([]byte("https://practicum.yandex.at/")))
	defer resp.Body.Close()

	for _, test := range testTable {
		resp, body := testRequest(t, ts, "GET", test.url, bytes.NewReader([]byte(test.body)))
		defer resp.Body.Close()
		assert.Equal(t, test.status, resp.StatusCode)
		if resp.StatusCode == 400 {
			assert.Equal(t, test.want, body)
		} else {
			assert.Equal(t, test.want, resp.Header.Get("Location"))
		}

	}
}

func TestShortenURL(t *testing.T) {
	t.Cleanup(func() {
		os.RemoveAll("test_storage2")
	})

	cfg := config.Config{
		BaseURL: "http://localhost:8080",
		StoragePath: "test_storage2",
	}

	log := logger.NewLogger()
	storage, _ := storage.NewStorage(&cfg)
	defer storage.File.Close()

	s := services.NewURLService(log, storage)

	tr := NewTransport(cfg, s, log)

	ts := httptest.NewServer(NewRouter(cfg, tr))
	defer ts.Close()

	var testTable = []struct {
		url    string
		body   string
		want   string
		status int
	}{
		{"/api/shorten", "https://practicum.yandex.ru/", "{\"result\":\"http://localhost:8080/5HZTcSScKLEKuTfRTUDifWUbO3kVswOjWFiZSB\"}", http.StatusCreated},
		{"/api/shorten", "https://practicum.yandex.at/", "{\"result\":\"http://localhost:8080/5HZTcSScKLEKuTfRTUDifWUbO3kVswOjWFdtYV\"}", http.StatusCreated},
		{"/api/shorten", "", "Empty body!", http.StatusBadRequest},
	}

	for _, test := range testTable {
		req := ShortenURLRequest{
			URL: test.body,
		}

		reqPayload, _ := json.Marshal(req)
		
		resp, body := testRequest(t, ts, "POST", test.url, bytes.NewReader(reqPayload))
		defer resp.Body.Close()
		assert.Equal(t, test.status, resp.StatusCode)
		assert.Equal(t, test.want, body)
	}
}

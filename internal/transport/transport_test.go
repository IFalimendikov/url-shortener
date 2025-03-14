package transport

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"

	// "strings"
	"testing"
	"url-shortener/internal/config"

	// "github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"url-shortener/internal/services"
	"url-shortener/internal/logger"
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
	cfg := config.Config{
		BaseURL: "http://localhost:8080",
	}

	log := logger.NewLogger()

	s := services.NewURLService(log)

	tr := NewTransport(cfg, s, log)

	ts := httptest.NewServer(NewRouter(cfg, tr))
	defer ts.Close()

	var testTable = []struct {
		url    string
		body   string
		want   string
		status int
	}{
		{"/", "https://practicum.yandex.ru/", "http://localhost:8080/1", http.StatusCreated},
		{"/", "https://practicum.yandex.at/", "http://localhost:8080/2", http.StatusCreated},
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
	cfg := config.Config{
		BaseURL: "http://localhost:8080",
	}

	log := logger.NewLogger()

	s := services.NewURLService(log)

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
		{"/1", "", "https://practicum.yandex.ru/", http.StatusTemporaryRedirect},
		{"/2", "", "https://practicum.yandex.at/", http.StatusTemporaryRedirect},
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

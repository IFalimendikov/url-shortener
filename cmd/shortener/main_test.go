package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	// "strings"
	"testing"

	// "github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader,) (*http.Response, string) {
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
    ts := httptest.NewServer(URLRouter())
    defer ts.Close()

    var testTable = []struct {
        url    string
        body   string
        want   string
        status int
    }{
        {"/", "https://practicum.yandex.ru/", "http://localhost:8080/4C92", http.StatusCreated},
        {"/", "https://practicum.yandex.at/", "http://localhost:8080/8OI4", http.StatusCreated},
        {"/", "", "\"Empty body!\"", http.StatusBadRequest},
        {"/", "practicum.yandex.ru/", "\"Mallformed URI!\"", http.StatusBadRequest},
    }

    for _, test := range testTable {
        resp, body := testRequest(t, ts, "POST", test.url, bytes.NewReader([]byte(test.body)))
        defer resp.Body.Close()
        assert.Equal(t, test.status, resp.StatusCode)
        assert.Equal(t, test.want, body)
    }
}

func TestGetURL(t *testing.T) {
    ts := httptest.NewServer(URLRouter())
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
        {"/4C92", "", "https://practicum.yandex.ru/", http.StatusTemporaryRedirect},
        {"/8OI4", "", "https://practicum.yandex.at/", http.StatusTemporaryRedirect},
        {"/91OP", "", "\"URL not found!\"", http.StatusBadRequest},
    }

    // First create the shortened URLs
    resp, _ := testRequest(t, ts, "POST", "/postURL", bytes.NewReader([]byte("https://practicum.yandex.ru/")))
    defer resp.Body.Close()
    resp, _ = testRequest(t, ts, "POST", "/postURL", bytes.NewReader([]byte("https://practicum.yandex.at/")))
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
package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	// "net/url"
	"strings"
	"testing"

	// "github.com/google/go-containerregistry/pkg/name"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostURL(t *testing.T) {
	type want struct {
		contentType string
		status      int
		body        string
	}

	tests := []struct {
		name    string
		request string
		body    string
		want    want
	}{
		{
			name:    "correct url #1",
			request: "/",
			body:    "https://practicum.yandex.ru/",
			want: want{
				contentType: "text/plain",
				status:      201,
				body:        "http://localhost:8080/4C92",
			},
		},
		{
			name:    "correct url #2",
			request: "/",
			body:    "https://practicum.yandex.at/",
			want: want{
				contentType: "text/plain",
				status:      201,
				body:        "http://localhost:8080/8OI4",
			},
		},
		{
			name:    "empty url",
			request: "/",
			body:    "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				status:      400,
				body:        "Empty URL!\n",
			},
		},
		{
			name:    "mallformed url",
			request: "/",
			body:    "practicum.yandex.ru/",
			want: want{
				contentType: "text/plain; charset=utf-8",
				status:      400,
				body:        "Mallformed URL!\n",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(test.body))

			w := httptest.NewRecorder()

			PostURL(w, req)

			result := w.Result()

			assert.Equal(t, test.want.status, result.StatusCode)
			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))

			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.body, string(body))
		})
	}
}

func TestGetURl(t *testing.T) {
	type want struct {
		status int
		body   string
		header string
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "correct test #1",
			request: "http://localhost:8080/4C92",
			want: want{
				status: 307,
				body:   "",
				header: "https://practicum.yandex.ru/",
			},
		},
		{
			name:    "correct test #2",
			request: "http://localhost:8080/8OI4",
			want: want{
				status: 307,
				body:   "",
				header: "https://practicum.yandex.at/",
			},
		},
		{
			name:    "incorrect test",
			request: "http://localhost:8080/91OP",
			want: want{
				status: 400,
				body:   "URL not found!\n",
				header: "",
			},
		},
		{
			name:    "empty url",
			request: "/",
			want: want{
				status: 400,
				body:   "URL is empty!\n",
				header: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.request, nil)

			w := httptest.NewRecorder()

			GetURL(w, req)

			result := w.Result()

			assert.Equal(t, test.want.status, result.StatusCode)
			assert.Equal(t, test.want.header, result.Header.Get("Location"))

			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.body, string(body))
		})
	}
}

// internal/handler/handler_test.go
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"url-shortener/internal/config"
	"url-shortener/internal/logger"
	"url-shortener/internal/models"
	"url-shortener/internal/services"
	"url-shortener/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) (*gin.Context, *httptest.ResponseRecorder, *Handler, config.Config) {
	t.Helper()

	cfg := config.Config{
		BaseURL:     "http://localhost:8080",
		StoragePath: "test_storage.json",
	}

	// Create a new recorder and context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Initialize actual components
	log := logger.NewLogger()
	store, err := storage.NewStorage(context.Background(), &cfg)
	require.NoError(t, err)

	service := services.NewURLService(context.Background(), log, store)
	h := NewHandler(service, log)

	return c, w, h, cfg
}

func TestNewHandler(t *testing.T) {
	_, _, h, _ := setupTest(t)
	assert.NotNil(t, h)
}

func TestPostURL(t *testing.T) {
	c, w, h, cfg := setupTest(t)

	// Test valid URL
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString("https://example.com"))
	c.Set("user_id", "test-user")
	h.PostURL(c, cfg)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Test empty body
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", nil)
	h.PostURL(c, cfg)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test invalid URL
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString("not-a-url"))
	h.PostURL(c, cfg)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Cleanup
	os.Remove(cfg.StoragePath)
}

func TestGetURL(t *testing.T) {
	c, w, h, cfg := setupTest(t)

	// First, save a URL
	userID := "test-user"

	// Create a proper request with context for POST
	postReq := httptest.NewRequest("POST", "/", bytes.NewBufferString("https://example.com"))
	postReq = postReq.WithContext(context.Background()) // Add this line
	c.Request = postReq
	c.Set("user_id", userID)
	h.PostURL(c, cfg)

	// Get the shortened URL from response
	shortURL := w.Body.String()
	shortID := shortURL[len(cfg.BaseURL)+1:]

	// Test getting the URL
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	// Create a proper request with context for GET
	getReq := httptest.NewRequest("GET", "/:id", nil)
	getReq = getReq.WithContext(context.Background()) // Add this line
	c.Request = getReq
	c.Params = []gin.Param{{Key: "id", Value: shortID}}
	h.GetURL(c)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, "https://example.com", c.Writer.Header().Get("Location"))

	// Test empty ID
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	emptyReq := httptest.NewRequest("GET", "/", nil)
	emptyReq = emptyReq.WithContext(context.Background()) // Add this line
	c.Request = emptyReq
	h.GetURL(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Cleanup
	os.Remove(cfg.StoragePath)
}

func TestShortenURL(t *testing.T) {
	c, w, h, cfg := setupTest(t)

	// Test valid request
	req := models.ShortenURLRequest{URL: "https://example.com"}
	body, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/api/shorten", bytes.NewBuffer(body))
	c.Set("user_id", "test-user")
	h.ShortenURL(c, cfg)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Test empty URL
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	req = models.ShortenURLRequest{URL: ""}
	body, _ = json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/api/shorten", bytes.NewBuffer(body))
	h.ShortenURL(c, cfg)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Cleanup
	os.Remove(cfg.StoragePath)
}

func TestShortenBatch(t *testing.T) {
	c, w, h, cfg := setupTest(t)

	// Test valid batch request
	req := []models.BatchUnitURLRequest{
		{ID: "1", URL: "https://example1.com"},
		{ID: "2", URL: "https://example2.com"},
	}
	body, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/api/shorten/batch", bytes.NewBuffer(body))
	c.Set("user_id", "test-user")
	h.ShortenBatch(c, cfg)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Test empty batch
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/shorten/batch", bytes.NewBuffer([]byte("[]")))
	h.ShortenBatch(c, cfg)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Cleanup
	os.Remove(cfg.StoragePath)
}

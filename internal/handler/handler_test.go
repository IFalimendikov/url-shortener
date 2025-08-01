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

	"github.com/brianvoe/gofakeit/v7"
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

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	log := logger.New()
	store, err := storage.New(context.Background(), &cfg)
	require.NoError(t, err)

	service := services.New(context.Background(), log, store)
	h := New(service, log)

	return c, w, h, cfg
}

func TestNewHandler(t *testing.T) {
	_, _, h, _ := setupTest(t)
	assert.NotNil(t, h)
}

func TestPostURL(t *testing.T) {
	c, w, h, cfg := setupTest(t)

	validURL := gofakeit.URL()
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString(validURL))
	c.Set("user_id", gofakeit.UUID())
	h.PostURL(c, cfg)
	assert.Equal(t, http.StatusCreated, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", nil)
	h.PostURL(c, cfg)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	invalidURL := gofakeit.Word()
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString(invalidURL))
	h.PostURL(c, cfg)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	os.Remove(cfg.StoragePath)
}

func TestGetURL(t *testing.T) {
	c, w, h, cfg := setupTest(t)

	userID := gofakeit.UUID()
	originalURL := gofakeit.URL()

	postReq := httptest.NewRequest("POST", "/", bytes.NewBufferString(originalURL))
	postReq = postReq.WithContext(context.Background())
	c.Request = postReq
	c.Set("user_id", userID)
	h.PostURL(c, cfg)

	shortURL := w.Body.String()
	shortID := shortURL[len(cfg.BaseURL)+1:]

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	getReq := httptest.NewRequest("GET", "/:id", nil)
	getReq = getReq.WithContext(context.Background())
	c.Request = getReq
	c.Params = []gin.Param{{Key: "id", Value: shortID}}
	h.GetURL(c)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, originalURL, c.Writer.Header().Get("Location"))

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	emptyReq := httptest.NewRequest("GET", "/", nil)
	emptyReq = emptyReq.WithContext(context.Background())
	c.Request = emptyReq
	h.GetURL(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	os.Remove(cfg.StoragePath)
}

func TestShortenURL(t *testing.T) {
	c, w, h, cfg := setupTest(t)

	req := models.ShortenURLRequest{URL: gofakeit.URL()}
	body, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/api/shorten", bytes.NewBuffer(body))
	c.Set("user_id", gofakeit.UUID())
	h.ShortenURL(c, cfg)
	assert.Equal(t, http.StatusCreated, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	req = models.ShortenURLRequest{URL: ""}
	body, _ = json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/api/shorten", bytes.NewBuffer(body))
	h.ShortenURL(c, cfg)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	os.Remove(cfg.StoragePath)
}

func TestShortenBatch(t *testing.T) {
	c, w, h, cfg := setupTest(t)

	req := []models.BatchUnitURLRequest{
		{ID: gofakeit.UUID(), URL: gofakeit.URL()},
		{ID: gofakeit.UUID(), URL: gofakeit.URL()},
	}
	body, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/api/shorten/batch", bytes.NewBuffer(body))
	c.Set("user_id", gofakeit.UUID())
	h.ShortenBatch(c, cfg)
	assert.Equal(t, http.StatusCreated, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/shorten/batch", bytes.NewBuffer([]byte("[]")))
	h.ShortenBatch(c, cfg)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	os.Remove(cfg.StoragePath)
}

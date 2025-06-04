package handler

import (
    "bytes"
    "encoding/json"
    "net/http/httptest"
    "testing"
    "url-shortener/internal/models"

    "github.com/brianvoe/gofakeit/v7"
    "github.com/gin-gonic/gin"
)

func BenchmarkPostURL(b *testing.B) {
    _, _, h, cfg := setupTest(&testing.T{})
    url := gofakeit.URL()
    userID := gofakeit.UUID()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        c.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString(url))
        c.Set("user_id", userID)
        h.PostURL(c, cfg)
    }
}

func BenchmarkGetURL(b *testing.B) {
    c, w, h, cfg := setupTest(&testing.T{})
    url := gofakeit.URL()
    userID := gofakeit.UUID()

    c.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString(url))
    c.Set("user_id", userID)
    h.PostURL(c, cfg)
    shortURL := w.Body.String()
    shortID := shortURL[len(cfg.BaseURL)+1:]

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        c.Request = httptest.NewRequest("GET", "/:id", nil)
        c.Params = []gin.Param{{Key: "id", Value: shortID}}
        h.GetURL(c)
    }
}

func BenchmarkShortenURL(b *testing.B) {
    _, _, h, cfg := setupTest(&testing.T{})
    req := models.ShortenURLRequest{URL: gofakeit.URL()}
    body, _ := json.Marshal(req)
    userID := gofakeit.UUID()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        c.Request = httptest.NewRequest("POST", "/api/shorten", bytes.NewBuffer(body))
        c.Set("user_id", userID)
        h.ShortenURL(c, cfg)
    }
}

func BenchmarkShortenBatch(b *testing.B) {
    _, _, h, cfg := setupTest(&testing.T{})
    req := []models.BatchUnitURLRequest{
        {ID: gofakeit.UUID(), URL: gofakeit.URL()},
        {ID: gofakeit.UUID(), URL: gofakeit.URL()},
    }
    body, _ := json.Marshal(req)
    userID := gofakeit.UUID()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        c.Request = httptest.NewRequest("POST", "/api/shorten/batch", bytes.NewBuffer(body))
        c.Set("user_id", userID)
        h.ShortenBatch(c, cfg)
    }
}
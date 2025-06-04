package handler

import (
    "bytes"
	"fmt"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/models"
)

func ExampleHandler_PostURL() {
    c, w, h, cfg := setupTest(&testing.T{})

    c.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString("https://github.com/IFalimendikov"))
    c.Set("user_id", "test-user")
    h.PostURL(c, cfg)
    
    fmt.Printf("Status: %d\n", w.Code)
    fmt.Printf("Response: %s\n", w.Body.String())
    
    // Output:
    // Status: 201
    // Response: http://localhost:8080/OlhrHnj0RPPBnBZbfX2w2whN8UX4ghviW6fXOasS2rO
}

func ExampleHandler_ShortenURL() {
    c, w, h, cfg := setupTest(&testing.T{})

    req := models.ShortenURLRequest{URL: "https://golang.org/doc"}
    body, _ := json.Marshal(req)
    c.Request = httptest.NewRequest("POST", "/api/shorten", bytes.NewBuffer(body))
    c.Set("user_id", "test-user")
    h.ShortenURL(c, cfg)

    fmt.Printf("Status: %d\n", w.Code)
    fmt.Printf("Response: %s\n", w.Body.String())

    // Output:
    // Status: 201
    // Response: {"result":"http://localhost:8080/4663gVcTVHSwM2qIUOBOCcjjTekACZ"}
}

func ExampleHandler_ShortenBatch() {
    c, w, h, cfg := setupTest(&testing.T{})

    req := []models.BatchUnitURLRequest{
        {ID: "1", URL: "https://ethereum.org"},
        {ID: "2", URL: "https://docs.soliditylang.org"},
    }
    body, _ := json.Marshal(req)
    c.Request = httptest.NewRequest("POST", "/api/shorten/batch", bytes.NewBuffer(body))
    c.Set("user_id", "test-user")
    h.ShortenBatch(c, cfg)

    fmt.Printf("Status: %d\n", w.Code)
    fmt.Printf("Response: %s\n", w.Body.String())

    // Output:
    // Status: 201
    // Response: [{"correlation_id":"1","short_url":"http://localhost:8080/Eu39w5QlXFlmYsZtRYfm2DHnBzD"},{"correlation_id":"2","short_url":"http://localhost:8080/LoYUInWAFzf5GRrJXD5h8KcAtZpxkzD2RfiL2m7"}]
}
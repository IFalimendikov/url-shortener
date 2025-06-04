package handler

import (
    "encoding/json"
    "errors"
    "io"
    "net/http"
    "url-shortener/internal/config"
    "url-shortener/internal/models"
    "url-shortener/internal/storage"

    "github.com/gin-gonic/gin"
)

// @Summary Shorten URL via JSON
// @Description Creates a shortened version of a URL provided in JSON format
// @Tags urls
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer JWT token"
// @Param request body models.ShortenURLRequest true "URL to shorten"
// @Success 201 {object} models.ShortenURLResponse "Shortened URL"
// @Success 409 {object} models.ShortenURLResponse "URL already exists"
// @Router /api/shorten [post]
func (t *Handler) ShortenURL(c *gin.Context, cfg config.Config) {
    var req models.ShortenURLRequest
    var res models.ShortenURLResponse

    body, err := io.ReadAll(c.Request.Body)
    if err != nil {
        c.String(http.StatusBadRequest, "Cant read body!")
        return
    }

    err = json.Unmarshal(body, &req)
    if err != nil {
        c.String(http.StatusBadRequest, "Couldn't unmarshal!")
        return
    }

    if req.URL == "" {
        c.String(http.StatusBadRequest, "Empty body!")
        return
    }

    userID := c.GetString("user_id")

    shortURL, err := t.service.SaveURL(c.Request.Context(), req.URL, userID)
    res.Result = cfg.BaseURL + "/" + string(shortURL)
    if err != nil {
        if errors.Is(err, storage.ErrorDuplicate) {
            c.JSON(http.StatusConflict, res)
            return
        }
        c.String(http.StatusBadRequest, "Couldn't encode URL!")
        return
    }

    c.JSON(http.StatusCreated, res)
}

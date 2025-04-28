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

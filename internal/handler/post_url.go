package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"url-shortener/internal/config"
	"url-shortener/internal/storage"

	"github.com/gin-gonic/gin"
)

func (t *Handler) PostURL(c *gin.Context, cfg config.Config) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Cant read body!")
		return
	}

	if len(body) == 0 {
		c.String(http.StatusBadRequest, "Empty body!")
		return
	}

	urlStr := string(body)
	parsedURL, err := url.Parse(urlStr)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		c.String(http.StatusBadRequest, "Malformed URI!")
		return
	}

	userID := c.GetString("user_id")

	shortURL, err := t.service.SaveURL(c.Request.Context(), urlStr, string(userID))
	shortURL = fmt.Sprintf("%s/%s", cfg.BaseURL, shortURL)
	if err != nil {
		if errors.Is(err, storage.ErrorDuplicate) {
			c.String(http.StatusConflict, shortURL)
			return
		}
		c.String(http.StatusBadRequest, "Couldn't encode URL!")
		return
	}
	c.String(http.StatusCreated, shortURL)
}

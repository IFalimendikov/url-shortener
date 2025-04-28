package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"url-shortener/internal/config"

	"github.com/gin-gonic/gin"
)

func (t *Handler) DeleteURLs(c *gin.Context, cfg config.Config) {
	var req []string

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Error reading body!")
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		c.String(http.StatusBadRequest, "Error unmarshalling body!")
		return
	}

	if len(req) == 0 {
		c.String(http.StatusBadRequest, "Empty or malformed body sent!")
		return
	}

	userID := c.GetString("user_id")

	go t.service.DeleteURLs(req, userID)

	c.Status(http.StatusAccepted)
}

package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"url-shortener/internal/config"

	"github.com/gin-gonic/gin"
)

// @Summary Delete URLs
// @Description Delete multiple URLs for a specific user
// @Tags urls
// @Accept json
// @Produce plain
// @Param Authorization header string true "Bearer JWT token"
// @Param request body []string true "Array of URLs to delete"
// @Success 202 {string} string "Accepted"
// @Failure 400 {string} string "Error reading body!/Error unmarshalling body!/Empty or malformed body sent!"
// @Router /api/urls [delete]
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

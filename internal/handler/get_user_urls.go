package handler

import (
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/models"

	"github.com/gin-gonic/gin"
)

// @Summary Get user's URLs
// @Description Retrieves all URLs associated with the authenticated user
// @Tags urls
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Success 200 {array} models.UserURLResponse "List of user's URLs"
// @Success 204 {string} string "No URLs found!"
// @Failure 400 {string} string "Error finding URLs!"
// @Router /api/user/urls [get]
func (t *Handler) GetUserURLs(c *gin.Context, cfg config.Config) {
	var res []models.UserURLResponse

	userID := c.GetString("user_id")

	err := t.service.GetUserURLs(c.Request.Context(), userID, &res)
	if err != nil {
		c.String(http.StatusBadRequest, "Error finding URLs!")
		return
	}

	if len(res) == 0 {
		c.String(http.StatusNoContent, "No URLs found!")
		return
	}

	for i := range res {
		res[i].ShortURL = cfg.BaseURL + "/" + res[i].ShortURL
	}

	c.JSON(http.StatusOK, res)
}

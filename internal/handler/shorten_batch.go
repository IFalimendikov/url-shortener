package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/models"

	"github.com/gin-gonic/gin"
)

func (t *Handler) ShortenBatch(c *gin.Context, cfg config.Config) {
	var req []models.BatchUnitURLRequest
	var res []models.BatchUnitURLResponse

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

	err = t.service.ShortenBatch(c.Request.Context(), userID, req, &res)
	if err != nil {
		c.String(http.StatusBadRequest, "Error saving URLs!")
		return
	}

	for i := range res {
		res[i].Short = cfg.BaseURL + "/" + res[i].Short
	}

	c.JSON(http.StatusCreated, res)
}

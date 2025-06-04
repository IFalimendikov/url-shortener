package handler

import (
    "encoding/json"
    "io"
    "net/http"
    "url-shortener/internal/config"
    "url-shortener/internal/models"

    "github.com/gin-gonic/gin"
)

// @Summary Shorten multiple URLs in batch
// @Description Creates shortened versions for multiple URLs in a single request
// @Tags urls
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer JWT token"
// @Param request body []models.BatchUnitURLRequest true "Array of URLs to shorten"
// @Success 201 {array} models.BatchUnitURLResponse "Array of shortened URLs"
// @Failure 400 {string} string "Error reading body!/Error unmarshalling body!/Empty or malformed body sent!/Error saving URLs!"
// @Router /api/shorten/batch [post]
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

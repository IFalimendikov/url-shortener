package handler

import (
    "errors"
    "net/http"
    "url-shortener/internal/storage"

    "github.com/gin-gonic/gin"
)

// @Summary Get original URL
// @Description Retrieves and redirects to the original URL from a shortened URL ID
// @Tags urls
// @Accept plain
// @Produce plain
// @Param id path string true "Shortened URL ID"
// @Success 307 {string} string "Temporary Redirect"
// @Failure 400 {string} string "URL not found!"
// @Failure 410 {string} string "URL was deleted!"
// @Header 307 {string} Location "Original URL for redirect"
// @Router /{id} [get]
func (t *Handler) GetURL(c *gin.Context) {
    id := c.Param("id")

    if id != "" {
        url, err := t.service.GetURL(c.Request.Context(), id)
        if err != nil {
            if errors.Is(err, storage.ErrorURLDeleted) {
                c.String(http.StatusGone, "URL was deleted!")
                return
            }
            c.String(http.StatusBadRequest, "URL not found!")
            return
        }

        c.Header("Location", url)
        c.Redirect(http.StatusTemporaryRedirect, url)

    } else {
        c.String(http.StatusBadRequest, "URL is empty!")
        return
    }
}
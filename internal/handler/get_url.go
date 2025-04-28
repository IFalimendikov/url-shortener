package handler

import (
	"errors"
	"net/http"
	"url-shortener/internal/storage"

	"github.com/gin-gonic/gin"
)

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

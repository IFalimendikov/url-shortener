package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (t *Handler) PingDB(c *gin.Context) {
	live := t.service.PingDB()
	if live {
		c.String(http.StatusOK, "Live")
		return
	}

	c.String(http.StatusInternalServerError, "Can't connect to the Database!")
}

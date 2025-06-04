package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// @Summary Ping database
// @Description Check if database connection is alive
// @Tags health
// @Accept plain
// @Produce plain
// @Success 200 {string} string "Live"
// @Failure 500 {string} string "Can't connect to the Database!"
// @Router /ping [get]
func (t *Handler) PingDB(c *gin.Context) {
    live := t.service.PingDB()
    if live {
        c.String(http.StatusOK, "Live")
        return
    }

    c.String(http.StatusInternalServerError, "Can't connect to the Database!")
}

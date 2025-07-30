package handler

import (
	"net"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/models"

	"github.com/gin-gonic/gin"
)

// @Summary Get service statistics
// @Description shows service stats
// @Tags stats
// @Accept plain
// @Produce plain
// @Success 200 {string} string "Sucess"
// @Failure 400 {string} string "Stats not found!"
// @Failure 403 {string} string "Forbidden IP!"
// @Router /api/internal/stats [get]
func (t *Handler) GetStats(c *gin.Context, cfg config.Config) {
	var res models.Stats

	if cfg.TrustedSubnet == "" {
		c.String(http.StatusForbidden, "No CIDR set!")
		return
	}

	s := c.GetHeader("X-Real-IP")
	userIP := net.ParseIP(s)

	_, network, err := net.ParseCIDR(cfg.TrustedSubnet)
	if err != nil {
		c.String(http.StatusInternalServerError, "Can't parse CIDR!")
		return
	}

	access := network.Contains(userIP)
	if !access {
		c.String(http.StatusForbidden, "IP address is not trusted!")
		return
	}

	stats, err := t.service.GetStats(c.Request.Context())
	if err != nil {
		c.String(http.StatusBadRequest, "Stats not found!")
		return
	}

	res.Urls = stats.Urls
	res.Users = stats.Users

	c.JSON(http.StatusOK, res)
}

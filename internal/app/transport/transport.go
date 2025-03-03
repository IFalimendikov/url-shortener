package transport

import (
	"io"
	"net/http"
	"strings"
	"url-shortener/internal/app/config"

	"github.com/gin-gonic/gin"
	base "github.com/jcoene/go-base62"
)

var urlMap = map[string]string{}
var counter int

func NewURLRouter(cfg config.Config) *gin.Engine {
	r := gin.Default()

	r.POST("/", func(c *gin.Context) {
		PostURL(c, cfg)
	})
	r.GET("/:id", GetURL)

	return r
}

func PostURL(c *gin.Context, cfg config.Config) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusBadRequest, "Only POST method allowed!")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Cant read body!")
		return
	}

	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, "Empty body!")
		return
	}

	urlStr := string(body)
	if !strings.HasPrefix(urlStr, "https://") && !strings.HasPrefix(urlStr, "http://") {
		c.JSON(http.StatusBadRequest, "Mallformed URI!")
		return
	}

	counter++

	urlShort := base.Encode(int64(counter * 1000000))
	urlMap[urlShort] = string(body)

	c.Header("Content-Type", "text/plain")
	c.String(http.StatusCreated, "%s/%s", cfg.BaseAddr, urlShort)
}

func GetURL(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusBadRequest, "Only GET method allowed!")
		return
	}

	id := c.Param("id")

	if id != "" {
		url, ok := urlMap[id]
		if !ok {
			c.JSON(http.StatusBadRequest, "URL not found!")
			return
		}

		c.Header("Location", url)
		c.Redirect(http.StatusTemporaryRedirect, url)

	} else {
		c.JSON(http.StatusBadRequest, "URL is empty!")
		return
	}
}

package transport

import (
	"io"
	"net/http"
	"sync"
	"net/url"
	"url-shortener/internal/app/config"

	"github.com/gin-gonic/gin"
	base "github.com/jcoene/go-base62"
)

type Storage struct {
	urls map[string]string
	counter int
	mu sync.RWMutex
}

func NewURLRouter(cfg config.Config) *gin.Engine {
	storage := &Storage{
		urls: make(map[string]string),
	}

	r := gin.Default()

	r.POST("/", func(c *gin.Context) {
		storage.PostURL(c, cfg)
	})
	r.GET("/:id", storage.GetURL)

	return r
}

func (s *Storage) PostURL(c *gin.Context, cfg config.Config) {
	if c.Request.Method != http.MethodPost {
		c.String(http.StatusBadRequest, "Only POST method allowed!")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Cant read body!")
		return
	}

	if len(body) == 0 {
		c.String(http.StatusBadRequest, "Empty body!")
		return
	}

	urlStr := string(body)
	parsedURL, err := url.Parse(urlStr)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		c.String(http.StatusBadRequest, "Malformed URI!")
		return
	}

	s.mu.Lock()
	s.counter++
	urlShort := base.Encode(int64(s.counter))
	s.urls[urlShort] = urlStr
	s.mu.Unlock()

	c.String(http.StatusCreated, "%s/%s", cfg.BaseURL, urlShort)
}

func (s *Storage) GetURL(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.String(http.StatusBadRequest, "Only GET method allowed!")
		return
	}

	id := c.Param("id")

	if id != "" {
		s.mu.RLock()
		url, ok := s.urls[id]
		s.mu.RUnlock()
		if !ok {
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

package transport

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
	"url-shortener/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type URLService interface {
	ServShort(url string) (string, error)
	ServSave(url string) (string, error)
	ServGet(shortURL string) (string, error)
}

type Transport struct {
	serviceURL URLService
	log        *zap.SugaredLogger
}

func NewTransport(cfg config.Config, s URLService, log *zap.SugaredLogger) Transport {
	return Transport{
		serviceURL: s,
		log:        log,
	}
}

func NewRouter(cfg config.Config, t Transport) *gin.Engine {

	r := gin.Default()

	r.Use(WithLogging(t.log))

	r.POST("/", func(c *gin.Context) {
		t.PostURL(c, cfg)
	})	

	r.POST("/api/shorten", func(c *gin.Context) {
		t.ShortenURL(c, cfg)
	})

	r.GET("/:id", t.GetURL)

	return r
}

func WithLogging(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		uri := c.Request.RequestURI
		method := c.Request.Method

		c.Next()

		status := c.Writer.Status()
		size := c.Writer.Size()

		latency := time.Since(start)
		log.Infoln(
			"uri", uri,
			"method", method,
			"duration", latency,
			"status", status,
			"size", size,
		)
	}
}

func (t *Transport) PostURL(c *gin.Context, cfg config.Config) {
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

	shortURL, err := t.serviceURL.ServSave(urlStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't encode URL!")
		return
	}

	c.String(http.StatusCreated, "%s/%s", cfg.BaseURL, shortURL)
}

func (t *Transport) GetURL(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.String(http.StatusBadRequest, "Only GET method allowed!")
		return
	}

	id := c.Param("id")

	if id != "" {
		url, err := t.serviceURL.ServGet(id)
		if err != nil {
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

func (t *Transport) ShortenURL(c *gin.Context, cfg config.Config) {
	var req *ShortenURLRequest
	var res ShortneURLResponse

	if c.Request.Method != http.MethodPost {
		c.String(http.StatusBadRequest, "Only POST method allowed!")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Cant read body!")
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't unmarshal!")
		return
	}

	if req.URL == "" {
		c.String(http.StatusBadRequest, "Empty body!")
		return
	}

	shortURL, err := t.serviceURL.ServShort(req.URL)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't encode URL!")
		return
	}

	res.Result = cfg.BaseURL + "/" + shortURL

	c.JSON(http.StatusCreated, res)
}

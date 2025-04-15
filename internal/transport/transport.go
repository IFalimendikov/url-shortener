package transport

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"log/slog"
	"net/http"
	"net/url"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/models"
	"url-shortener/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type URLService interface {
	ServSave(ctx context.Context, url, userID string) (string, error)
	ServGet(shortURL string) (string, error)
	ShortenBatch(ctx context.Context, req []models.BatchUnitURLRequest, res *[]models.BatchUnitURLResponse) error
	GetUserURLs(ctx context.Context, userID string, res *[]models.UserURLResponse) error
	PingDB() bool
}

type Transport struct {
	serviceURL URLService
	log        *slog.Logger
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

type gzipWriter struct {
	gin.ResponseWriter
	gzip *gzip.Writer
}

func NewTransport(cfg config.Config, s URLService, log *slog.Logger) Transport {
	return Transport{
		serviceURL: s,
		log:        log,
	}
}

func NewRouter(cfg config.Config, t Transport) *gin.Engine {

	r := gin.Default()

	r.Use(WithLogging(t.log))
	r.Use(WithDecodingReq())
	r.Use(WithEncodingRes())
	// r.Use(WithCookies())

	r.POST("/", func(c *gin.Context) {
		t.PostURL(c, cfg)
	})
	r.POST("/api/shorten", func(c *gin.Context) {
		t.ShortenURL(c, cfg)
	})
	r.POST("/api/shorten/batch", func(c *gin.Context) {
		t.ShortenBatch(c, cfg)
	})

	r.GET("/:id", t.GetURL)
	r.GET("/ping", t.PingDB)
	r.GET("/api/user/urls", t.GetUserURLs)

	return r
}

func WithLogging(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		uri := c.Request.RequestURI
		method := c.Request.Method

		c.Next()

		status := c.Writer.Status()
		size := c.Writer.Size()

		latency := time.Since(start)
        log.Info("request completed",
            "uri", uri,
            "method", method,
            "duration", latency.String(),
            "status", status,
            "size", size,
        )
		c.Next()
	}
}

func WithDecodingReq() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Content-Encoding")
		if header != "gzip" {
			c.Next()
			return
		}

		r, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			c.Next()
			return
		}
		defer r.Close()

		newBody, err := io.ReadAll(r)
		if err != nil {
			c.Next()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(newBody))
		c.Next()
	}
}

func WithEncodingRes() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Accept-Encoding")
		if header != "gzip" {
			c.Next()
			return
		}

		header = c.Writer.Header().Get("Content-Type")
		if header == "application/json" || header == "text/html" {
			gz := gzip.NewWriter(c.Writer)
			defer gz.Close()

			c.Header("Content-Encoding", "gzip")

			c.Writer = gzipWriter{
				ResponseWriter: c.Writer,
				gzip:           gz,
			}
		}
		c.Next()
	}
}

func WithCookies() gin.HandlerFunc {
	return func(c *gin.Context) {
		var UserID int
		if cookie, err := c.Cookie("jwt"); err == nil {
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(cookie, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					c.String(http.StatusBadRequest, "Unexpected signing method")
					return nil, err
				}
				return []byte("123"), nil
			})

			if err != nil {
				if claims.UserID == 0 {
					c.String(http.StatusUnauthorized, "User ID not found!")
					return
				}
			} else if token.Valid {
				UserID = claims.UserID
				c.Set("user_id", claims.UserID)
				c.Next()
			}
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
			},
			UserID: UserID,
		})

		signedToken, err := token.SignedString([]byte("123"))
		if err != nil {
			c.String(http.StatusInternalServerError, "Error signing token")
			return
		}

		c.Set("user_id", UserID)
		c.SetCookie("jwt", signedToken, 60, "/", "", false, true)
		c.Next()
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

	userID := c.GetHeader("user_id")

	shortURL, err := t.serviceURL.ServSave(c.Request.Context(), urlStr, userID)
	shortURL = fmt.Sprintf("%s/%s", cfg.BaseURL, shortURL)
	if err != nil {
		if errors.Is(err, services.ErrorDuplicate) {
			c.String(http.StatusConflict, shortURL)
			return
		}
		c.String(http.StatusBadRequest, "Couldn't encode URL!")
		return
	}
	c.String(http.StatusCreated, shortURL)
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
	var req models.ShortenURLRequest
	var res models.ShortenURLResponse

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

	userID := c.GetHeader("user_id")

	shortURL, err := t.serviceURL.ServSave(c.Request.Context(), req.URL, userID)
	res.Result = cfg.BaseURL + "/" + string(shortURL)
	if err != nil {
		if errors.Is(err, services.ErrorDuplicate) {
			c.JSON(http.StatusConflict, res)
			return
		}
		c.String(http.StatusBadRequest, "Couldn't encode URL!")
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (t *Transport) PingDB(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.String(http.StatusBadRequest, "Only GET method allowed!")
		return
	}

	live := t.serviceURL.PingDB()
	if live {
		c.String(http.StatusOK, "Live")
		return
	}

	c.String(http.StatusInternalServerError, "Can't connect to the Database!")
}

func (t *Transport) ShortenBatch(c *gin.Context, cfg config.Config) {
	var req []models.BatchUnitURLRequest
	var res []models.BatchUnitURLResponse

	if c.Request.Method != http.MethodPost {
		c.String(http.StatusBadRequest, "Only POST method allowed!")
		return
	}

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
		c.String(http.StatusBadRequest, "Empty or mallformed body sent!")
		return
	}

	err = t.serviceURL.ShortenBatch(c.Request.Context(), req, &res)
	if err != nil {
		c.String(http.StatusBadRequest, "Error saving URLs!")
		return
	}

	for i := range res {
		res[i].Short = cfg.BaseURL + "/" + res[i].Short
	}

	c.JSON(http.StatusCreated, res)
}

func (t *Transport) GetUserURLs(c *gin.Context) {
	var res []models.UserURLResponse

	if c.Request.Method != http.MethodGet {
		c.String(http.StatusBadRequest, "Only GET method allowed!")
		return
	}

	userID := c.GetHeader("user_id")

	err := t.serviceURL.GetUserURLs(c.Request.Context(), userID, &res)
	if err != nil {
		c.String(http.StatusBadRequest, "Error finding URLs!")
		return
	}


	c.JSON(http.StatusCreated, res)
}
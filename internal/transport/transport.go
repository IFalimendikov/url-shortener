package transport

import (
	"bytes"
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Transport handles HTTP transport layer operations including middleware and routing
type Transport struct {
	handler *handler.Handler
	log     *slog.Logger
	cfg     config.Config
}

// Claims represents JWT claims structure with user identification
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// gzipWriter wraps gin.ResponseWriter to provide gzip compression
type gzipWriter struct {
	gin.ResponseWriter
	gzip *gzip.Writer
}

// Write implements io.Writer interface for gzipWriter
func (gz gzipWriter) Write(data []byte) (int, error) {
	return gz.gzip.Write(data)
}

// New creates a new Transport instance with the provided configuration and handlers
func New(cfg config.Config, h *handler.Handler, log *slog.Logger) *Transport {
	return &Transport{
		handler: h,
		log:     log,
		cfg:     cfg,
	}
}

// NewRouter sets up and configures a new gin router with all necessary middlewares and routes
func NewRouter(t *Transport) *gin.Engine {

	r := gin.Default()

	r.Use(t.WithLogging(t.log))
	r.Use(t.WithDecodingReq())
	r.Use(t.WithEncodingRes())
	r.Use(t.WithCookies())

	r.POST("/", func(c *gin.Context) {
		t.handler.PostURL(c, t.cfg)
	})
	r.POST("/api/shorten", func(c *gin.Context) {
		t.handler.ShortenURL(c, t.cfg)
	})
	r.POST("/api/shorten/batch", func(c *gin.Context) {
		t.handler.ShortenBatch(c, t.cfg)
	})

	r.GET("/:id", t.handler.GetURL)
	r.GET("/ping", t.handler.PingDB)
	r.GET("/api/user/urls", func(c *gin.Context) {
		t.handler.GetUserURLs(c, t.cfg)
	})
	r.GET("/api/internal/stats", func(c *gin.Context) {
		t.handler.GetStats(c, t.cfg)
	} )

	r.DELETE("/api/user/urls", func(c *gin.Context) {
		t.handler.DeleteURLs(c, t.cfg)
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

// WithLogging adds request logging middleware that records URI, method, duration, status, and size.
func (t *Transport) WithLogging(log *slog.Logger) gin.HandlerFunc {
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

// WithDecodingReq adds middleware to handle gzip-encoded request bodies.
func (t *Transport) WithDecodingReq() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Content-Encoding")
		if header != "gzip" {
			c.Next()
			return
		}

		r, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			slog.Error("failed to create new gzip reader",
				"error", err,
				"path", c.Request.URL.Path)
			c.Next()
			return
		}
		defer r.Close()

		newBody, err := io.ReadAll(r)
		if err != nil {
			slog.Error("failed to read gzipped body",
				"error", err,
				"path", c.Request.URL.Path)
			c.Next()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(newBody))
		c.Next()
	}
}

// WithEncodingRes adds middleware to handle gzip encoding of response bodies for JSON and HTML content.
func (t *Transport) WithEncodingRes() gin.HandlerFunc {
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

// WithCookies adds middleware to handle JWT authentication via cookies and user identification.
func (t *Transport) WithCookies() gin.HandlerFunc {
	return func(c *gin.Context) {
		var UserID string
		if cookie, err := c.Cookie("jwt"); err == nil {
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(cookie, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					c.String(http.StatusBadRequest, "Unexpected signing method!")
					return nil, err
				}
				return []byte("123"), nil
			})

			if err != nil {
				if claims.UserID == "" {
					c.String(http.StatusUnauthorized, "User ID not found!")
					return
				}
			} else if token.Valid {
				UserID = claims.UserID
				c.Set("user_id", UserID)
				c.Next()
				return
			}
		}

		UserID = uuid.NewString()

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
			},
			UserID: UserID,
		})

		signedToken, err := token.SignedString([]byte("123"))
		if err != nil {
			slog.Error("failed to sign token",
				"error", err,
				"path", c.Request.URL.Path)
			c.Next()
			return
		}

		c.Set("user_id", UserID)
		c.SetCookie("jwt", signedToken, 60, "/", "", false, true)
		c.Next()
	}
}

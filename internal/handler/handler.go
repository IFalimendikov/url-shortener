package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"url-shortener/internal/config"
	"url-shortener/internal/models"
	"url-shortener/internal/storage"

	"github.com/gin-gonic/gin"
)

type URLService interface {
	SaveURL(ctx context.Context, url, userID string) (string, error)
	GetURL(ctx context.Context, shortURL string) (string, error)
	ShortenBatch(ctx context.Context, userID string, req []models.BatchUnitURLRequest, res *[]models.BatchUnitURLResponse) error
	GetUserURLs(ctx context.Context, userID string, res *[]models.UserURLResponse) error
	PingDB() bool
	DeleteURLs(req []string, userID string) error
}

type Handler struct {
	serviceURL URLService
	log        *slog.Logger
}

func NewHandler(s URLService, log *slog.Logger) *Handler {
	return &Handler{
		serviceURL: s,
		log:        log,
	}
}

func (t *Handler) PostURL(c *gin.Context, cfg config.Config) {
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

	userID := c.GetString("user_id")

	shortURL, err := t.serviceURL.SaveURL(c.Request.Context(), urlStr, string(userID))
	shortURL = fmt.Sprintf("%s/%s", cfg.BaseURL, shortURL)
	if err != nil {
		if errors.Is(err, storage.ErrorDuplicate) {
			c.String(http.StatusConflict, shortURL)
			return
		}
		c.String(http.StatusBadRequest, "Couldn't encode URL!")
		return
	}
	c.String(http.StatusCreated, shortURL)
}

func (t *Handler) GetURL(c *gin.Context) {
	id := c.Param("id")

	if id != "" {
		url, err := t.serviceURL.GetURL(c.Request.Context(), id)
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

func (t *Handler) ShortenURL(c *gin.Context, cfg config.Config) {
	var req models.ShortenURLRequest
	var res models.ShortenURLResponse

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

	userID := c.GetString("user_id")

	shortURL, err := t.serviceURL.SaveURL(c.Request.Context(), req.URL, userID)
	res.Result = cfg.BaseURL + "/" + string(shortURL)
	if err != nil {
		if errors.Is(err, storage.ErrorDuplicate) {
			c.JSON(http.StatusConflict, res)
			return
		}
		c.String(http.StatusBadRequest, "Couldn't encode URL!")
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (t *Handler) PingDB(c *gin.Context) {
	live := t.serviceURL.PingDB()
	if live {
		c.String(http.StatusOK, "Live")
		return
	}

	c.String(http.StatusInternalServerError, "Can't connect to the Database!")
}

func (t *Handler) ShortenBatch(c *gin.Context, cfg config.Config) {
	var req []models.BatchUnitURLRequest
	var res []models.BatchUnitURLResponse

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
		c.String(http.StatusBadRequest, "Empty or malformed body sent!")
		return
	}

	userID := c.GetString("user_id")

	err = t.serviceURL.ShortenBatch(c.Request.Context(), userID, req, &res)
	if err != nil {
		c.String(http.StatusBadRequest, "Error saving URLs!")
		return
	}

	for i := range res {
		res[i].Short = cfg.BaseURL + "/" + res[i].Short
	}

	c.JSON(http.StatusCreated, res)
}

func (t *Handler) GetUserURLs(c *gin.Context, cfg config.Config) {
	var res []models.UserURLResponse

	userID := c.GetString("user_id")

	err := t.serviceURL.GetUserURLs(c.Request.Context(), userID, &res)
	if err != nil {
		c.String(http.StatusBadRequest, "Error finding URLs!")
		return
	}

	if len(res) == 0 {
		c.String(http.StatusNoContent, "No URLs found!")
		return
	}

	for i := range res {
		res[i].ShortURL = cfg.BaseURL + "/" + res[i].ShortURL
	}

	c.JSON(http.StatusOK, res)
}

func (t *Handler) DeleteURLs(c *gin.Context, cfg config.Config) {
	var req []string

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
		c.String(http.StatusBadRequest, "Empty or malformed body sent!")
		return
	}

	userID := c.GetString("user_id")

	go t.serviceURL.DeleteURLs(req, userID)

	c.Status(http.StatusAccepted)
}

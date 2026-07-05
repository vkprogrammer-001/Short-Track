package handler

import (
    "net/http"
    "short-track/internal/domain"
    "github.com/gin-gonic/gin"
)

type Handler struct {
    svc domain.ShortenerService
}

func NewHandler(svc domain.ShortenerService) *Handler {
    return &Handler{svc: svc}
}

func (h *Handler) ShortenURL(c *gin.Context) {
    var req struct {
        URL string `json:"url" binding:"required,url"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
        return
    }

    code, err := h.svc.Shorten(req.URL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"short_url": "http://localhost:8080/" + code})
}

func (h *Handler) Redirect(c *gin.Context) {
    code := c.Param("code")
    originalURL, err := h.svc.Resolve(code)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    c.Redirect(http.StatusMovedPermanently, originalURL)
}
package handler

import (
	"net/http"

	"zoomgate/internal/store"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	store store.Store
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(s store.Store) *HealthHandler {
	return &HealthHandler{store: s}
}

func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *HealthHandler) Readiness(c *gin.Context) {
	if s, ok := h.store.(*store.SQLiteStore); ok {
		db, err := s.DB().DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "detail": "db connection failed"})
			return
		}
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "detail": "db ping failed"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

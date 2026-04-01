package handler

import (
	"net/http"

	"zoomgate/internal/store"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles liveness and readiness health check requests.
type HealthHandler struct {
	dataStore store.Store
}

// NewHealthHandler creates a new HealthHandler with the given data store.
func NewHealthHandler(dataStore store.Store) *HealthHandler {
	return &HealthHandler{dataStore: dataStore}
}

// Liveness returns HTTP 200 if the server process is alive.
func (handler *HealthHandler) Liveness(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Readiness returns HTTP 200 if the server is ready to accept traffic,
// including verifying database connectivity.
func (handler *HealthHandler) Readiness(ctx *gin.Context) {
	if sqliteStore, ok := handler.dataStore.(*store.SQLiteStore); ok {
		sqlDB, err := sqliteStore.GormDB().DB()
		if err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "detail": "db connection failed"})
			return
		}
		if err := sqlDB.Ping(); err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "detail": "db ping failed"})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ready"})
}

package server

import (
	"zoomgate/internal/handler"
)

func (s *Server) setupRoutes() {
	// Health check endpoints - no auth required
	healthCheck := handler.NewHealthHandler(s.store)
	s.engine.GET("/healthz", healthCheck.Liveness)
	s.engine.GET("/readyz", healthCheck.Readiness)
}

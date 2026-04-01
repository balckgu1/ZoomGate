package server

import (
	"zoomgate/internal/handler"
	"zoomgate/internal/middleware"
	"zoomgate/internal/model"
)

// setupRoutes registers all HTTP routes and attaches middleware chains.
func (srv *Server) setupRoutes() {
	// Global middleware applied to all routes
	srv.engine.Use(middleware.RequestID())
	srv.engine.Use(middleware.CORS())

	// Health check endpoints — no auth required
	healthHandler := handler.NewHealthHandler(srv.dataStore)
	srv.engine.GET("/healthz", healthHandler.Liveness)
	srv.engine.GET("/readyz", healthHandler.Readiness)

	// Admin login endpoint — no auth required
	adminAuthHandler := handler.NewAdminAuthHandler(srv.dataStore, srv.appConfig, srv.logger)
	srv.engine.POST("/admin/api/auth/login", adminAuthHandler.Login)

	// Admin API group — requires JWT auth with admin role
	adminGroup := srv.engine.Group("/admin/api")
	adminGroup.Use(middleware.JWTAuth(srv.appConfig.Auth.JWTSecret, srv.dataStore))
	adminGroup.Use(middleware.RequireRole(string(model.RoleAdmin)))
	{
		// User management
		userHandler := handler.NewUserHandler(srv.dataStore, srv.appConfig, srv.logger)
		adminGroup.GET("/users", userHandler.List)
		adminGroup.POST("/users", userHandler.Create)
		adminGroup.PUT("/users/:id", userHandler.Update)
		adminGroup.DELETE("/users/:id", userHandler.Delete)
		adminGroup.POST("/users/:id/api-key", userHandler.GenerateAPIKey)
	}
}

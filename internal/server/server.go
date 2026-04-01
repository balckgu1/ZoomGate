package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"zoomgate/internal/config"
	"zoomgate/internal/store"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Server wraps the HTTP server, Gin engine, data store, and logger together.
type Server struct {
	appConfig  *config.Config
	engine     *gin.Engine
	httpServer *http.Server
	dataStore  store.Store
	logger     *zap.Logger
}

// New creates and configures a new Server instance with the given
// configuration, data store, and logger. It sets up routing and middleware.
func New(appConfig *config.Config, dataStore store.Store, logger *zap.Logger) *Server {
	// Set Gin mode based on configuration
	if appConfig.Logging.Level == "debug" || appConfig.Server.DevMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	// Recover from panics and return 500
	engine.Use(gin.Recovery())

	srv := &Server{
		appConfig: appConfig,
		engine:    engine,
		dataStore: dataStore,
		logger:    logger,
	}

	// Register all routes and middleware chains
	srv.setupRoutes()

	srv.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", appConfig.Server.Port),
		Handler:      engine,
		ReadTimeout:  appConfig.Server.ReadTimeout,
		WriteTimeout: appConfig.Server.WriteTimeout,
	}

	return srv
}

// Start begins listening and serving HTTP requests. It blocks until
// the server is shut down or encounters a fatal error.
func (srv *Server) Start() error {
	srv.logger.Info("starting server", zap.Int("port", srv.appConfig.Server.Port))
	err := srv.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server listen: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the server with a 10-second timeout,
// allowing in-flight requests to complete.
func (srv *Server) Shutdown() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.logger.Info("shutting down server...")
	return srv.httpServer.Shutdown(shutdownCtx)
}

// Engine returns the underlying Gin engine for testing or advanced configuration.
func (srv *Server) Engine() *gin.Engine {
	return srv.engine
}

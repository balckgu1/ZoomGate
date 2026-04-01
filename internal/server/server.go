package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"zoomgate/internal/config"
	"zoomgate/internal/middleware"
	"zoomgate/internal/store"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	cfg    *config.Config
	engine *gin.Engine
	http   *http.Server
	store  store.Store
	logger *zap.Logger
}

// New is a constructor for a new Server instance.
func New(cfg *config.Config, st store.Store, logger *zap.Logger) *Server {
	// Set Gin mode based on configuration
	if cfg.Logging.Level == "debug" || cfg.Server.DevMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin engine
	engine := gin.New()
	// Add middleware to the engine
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestID())
	engine.Use(middleware.CORS())

	s := &Server{
		cfg:    cfg,
		engine: engine,
		store:  st,
		logger: logger,
	}

	// Add route
	s.setupRoutes()

	s.http = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return s
}

// Start the Http server.
func (s *Server) Start() error {
	s.logger.Info("starting server", zap.Int("port", s.cfg.Server.Port))
	err := s.http.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server listen: %w", err)
	}
	return nil
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.logger.Info("shutting down server...")
	return s.http.Shutdown(ctx)
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}

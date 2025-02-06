package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/k6zma/DockerMonitoringApp/backend/internal/application/usecases"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/infrastructure/config"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/infrastructure/db/postgres/repositories"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/presentation/handlers"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/presentation/routes"
	"github.com/k6zma/DockerMonitoringApp/backend/pkg/utils"
)

type Server struct {
	httpServer *http.Server
	logger     utils.LoggerInterface
}

func NewServer(cfg *config.Config, db *sqlx.DB, logger utils.LoggerInterface) *Server {
	repo := repositories.NewContainerStatusRepositoryImpl(db, logger)
	useCase := usecases.NewContainerStatusUseCase(repo, logger)
	containerHandler := handlers.NewContainerStatusHandler(useCase, logger)
	errHandler := handlers.NewErrorHandlers(logger)

	router := routes.InitRoutes(cfg, errHandler, containerHandler, logger)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		logger:     logger,
	}
}

func (s *Server) Start() error {
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Infof("SERVER: failed to start HTTP server: %v\n", err)
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

func (s *Server) Stop() error {
	if err := s.httpServer.Close(); err != nil {
		s.logger.Infof("SERVER: failed to stop HTTP server: %v\n", err)
		return fmt.Errorf("failed to stop HTTP server: %w", err)
	}

	return nil
}

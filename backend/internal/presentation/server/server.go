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
}

func NewServer(cfg *config.Config, db *sqlx.DB) *Server {
	repo := repositories.NewContainerStatusRepositoryImpl(db)
	useCase := usecases.NewContainerStatusUseCase(repo)
	handler := handlers.NewContainerStatusHandler(useCase)

	router := routes.InitRoutes(cfg, handler)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
	}
}

func (s *Server) Start() error {
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		utils.LoggerInstance.Infof("SERVER: failed to start HTTP server: %v\n", err)
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

func (s *Server) Stop() error {
	if err := s.httpServer.Close(); err != nil {
		utils.LoggerInstance.Infof("SERVER: failed to stop HTTP server: %v\n", err)
		return fmt.Errorf("failed to stop HTTP server: %w", err)
	}

	return nil
}

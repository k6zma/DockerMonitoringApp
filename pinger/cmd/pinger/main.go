package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/k6zma/DockerMonitoringApp/pinger/internal/application/usecases"
	"github.com/k6zma/DockerMonitoringApp/pinger/internal/infrastructure/backend"
	"github.com/k6zma/DockerMonitoringApp/pinger/internal/infrastructure/config"
	"github.com/k6zma/DockerMonitoringApp/pinger/internal/infrastructure/docker"
	"github.com/k6zma/DockerMonitoringApp/pinger/internal/infrastructure/flags"
	"github.com/k6zma/DockerMonitoringApp/pinger/pkg/utils"
)

func main() {
	flagsData, err := flags.ParseFlags()
	if err != nil {
		panic(fmt.Errorf("flags parsing failed: %w", err))
	}

	logger, err := utils.NewLogger(flagsData.LoggerLevel)
	if err != nil {
		panic(fmt.Errorf("logger init failed: %w", err))
	}

	logger.Infof("Loading config from %s", flagsData.ConfigFilePath)
	cfg, err := config.Load(flagsData.ConfigFilePath)
	if err != nil {
		logger.Fatalf("Config error: %v", err)
	}
	logger.Infof("Config loaded: Backend - %+v, Ping - %+v, Docker - %+v", *cfg.Backend, *cfg.Ping, *cfg.Docker)

	containerRepo, err := docker.NewDockerContainerRepo(cfg, logger)
	if err != nil {
		logger.Fatalf("Docker repository init failed: %v", err)
	}

	statusRepo := backend.NewBackendStatusRepo(
		cfg.Backend.URL,
		cfg.Backend.APIKey,
		logger,
	)

	pinger := usecases.NewPingerUsecase(
		containerRepo,
		statusRepo,
		cfg.Ping.PingInterval,
		logger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-stop
		logger.Infof("Received signal: %v", sig)
		cancel()
	}()

	logger.Info("Starting pinger service")
	if err := pinger.Run(ctx); err != nil {
		logger.Fatalf("Pinger service failed: %v", err)
	}
}

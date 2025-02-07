package main

import (
	"context"
	"log"
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
	flags, err := flags.ParseFlags()
	if err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}

	logger, err := utils.NewLogger(flags.LoggerLevel)
	if err != nil {
		log.Fatalf("Logger initialization failed: %v", err)
	}

	cfg, err := config.Load(flags.ConfigFilePath)
	if err != nil {
		logger.Fatalf("Config error: %v", err)
	}

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
		cfg.PingInterval,
		logger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupShutdownHandler(cancel, logger)

	if err := pinger.Run(ctx); err != nil {
		logger.Fatalf("Pinger service failed: %v", err)
	}
}

func setupShutdownHandler(cancel context.CancelFunc, logger utils.LoggerInterface) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Infof("Received signal: %v", sig)
		cancel()
	}()
}

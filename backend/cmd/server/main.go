package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/k6zma/DockerMonitoringApp/backend/docs"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/infrastructure/config"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/infrastructure/db/postgres"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/infrastructure/flags"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/infrastructure/migrations"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/presentation/server"
	"github.com/k6zma/DockerMonitoringApp/backend/pkg/utils"
)

// @title Docker Monitoring API
// @version 1.0
// @description REST API for monitoring Docker containers.
// @contact.name Mikhail Gunin
// @contact.url https://github.com/k6zma
// @contact.email k6zma@yandex.ru
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Api-Key
// @security ApiKeyAuth
func main() {
	appFlags, err := flags.ParseFlags()
	if err != nil {
		panic(err)
	}

	logger, err := utils.NewLogger(appFlags.LoggerLevel)
	if err != nil {
		panic(err)
	}

	utils.LoggerInstance.Infof("ENTRY POINT: parsed flags: Config path - %+v, Logging level - %+v",
		appFlags.ConfigFilePath, appFlags.LoggerLevel)
	cfg, err := config.LoadConfig(appFlags.ConfigFilePath)
	if err != nil {
		utils.LoggerInstance.Fatalf("ENTRY POINT: failed to load configuration: %v", err)
	}
	utils.LoggerInstance.Infof("ENTRY POINT: loaded configuration: Server - %+v, DB - %+v, MigrationsConfig - %+v, API Key - %+v",
		cfg.Server, cfg.DB, cfg.MigrationsConfig, cfg.AuthAPI)

	utils.LoggerInstance.Infof("ENTRY POINT: conecting to database \"%s:%d\"", cfg.DB.Host, cfg.DB.Port)
	database, err := postgres.NewPsqlDB(cfg.DB)
	if err != nil {
		utils.LoggerInstance.Fatalf("ENTRY POINT: failed to connect to database: %v", err)
	}

	defer func() {
		utils.LoggerInstance.Info("ENTRY POINT: closing database connection")

		err = database.Close()
		if err != nil {
			utils.LoggerInstance.Fatalf("ENTRY POINT: failed to close database connection: %v", err)
		}
	}()
	utils.LoggerInstance.Info("ENTRY POINT: database connected successfully")

	utils.LoggerInstance.Infof("ENTRY POINT: applying migrations from \"%s\" folder", cfg.MigrationsConfig.Path)

	mig := migrations.NewMigrate(database, cfg.MigrationsConfig.Path, logger)
	switch cfg.MigrationsConfig.Type {
	case "apply":
		err = mig.ApplyMigrations()
	case "drop":
		err = mig.DropMigrations()
	case "rollback":
		err = mig.RollbackMigrations()
	default:
		utils.LoggerInstance.Fatalf("ENTRY POINT: unknown migrations type: %s", cfg.MigrationsConfig.Path)
	}

	if err != nil {
		utils.LoggerInstance.Fatalf("ENTRY POINT: failed to apply migrations: %v", err)
	}
	utils.LoggerInstance.Info("ENTRY POINT: migrations applied successfully")

	utils.LoggerInstance.Infof("ENTRY POINT: starting server on port \"localhost:%d\"", cfg.Server.Port)
	serv := server.NewServer(cfg, database, logger)
	go func() {
		if err := serv.Start(); err != nil {
			utils.LoggerInstance.Fatalf("ENTRY POINT: failed to start server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	utils.LoggerInstance.Info("ENTRY POINT: shutting down server")

	if err := serv.Stop(); err != nil {
		utils.LoggerInstance.Fatalf("ENTRY POINT: failed to stop server: %v", err)
	}
	utils.LoggerInstance.Info("ENTRY POINT: server stopped successfully")
}

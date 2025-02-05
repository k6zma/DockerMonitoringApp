package flags

import (
	"flag"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type AppFlags struct {
	ConfigFilePath string `validate:"required,file"`
	MigrationsPath string `validate:"required,dir"`
	MigrationsType string `validate:"oneof=apply drop rollback"`
	LoggerLevel    string `validate:"oneof=debug info warn error dpanic panic fatal"`
}

func ParseFlags() (*AppFlags, error) {
	configFile := flag.String("config_path", "config.json", "Path to json configuration file")
	migrationsPath := flag.String("migrations_path", "migrations", "Path to migrations directory")
	migrationsType := flag.String("migrations_type", "file", "Type of migrations (apply, drop, rollback)")
	loggerLevel := flag.String("logger_level", "debug", "Logger level (debug, info, warn, error, dpanic, panic, fatal)")

	flag.Parse()

	appFlags := &AppFlags{
		ConfigFilePath: *configFile,
		MigrationsPath: *migrationsPath,
		MigrationsType: *migrationsType,
		LoggerLevel:    *loggerLevel,
	}

	validate := validator.New()
	if err := validate.Struct(appFlags); err != nil {
		return nil, fmt.Errorf("invalid flags: %w", err)
	}

	return appFlags, nil
}

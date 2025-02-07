package flags

import (
	"flag"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type AppFlags struct {
	ConfigFilePath string `validate:"required,file"`
	LoggerLevel    string `validate:"oneof=debug info warn error dpanic panic fatal"`
}

func ParseFlags() (*AppFlags, error) {
	configFile := flag.String("config_path", "config.json", "Path to config file")
	loggerLevel := flag.String("logger_level", "debug", "Logging level")

	flag.Parse()

	flags := &AppFlags{
		ConfigFilePath: *configFile,
		LoggerLevel:    *loggerLevel,
	}

	validate := validator.New()
	if err := validate.Struct(flags); err != nil {
		return nil, fmt.Errorf("flags validation error: %w", err)
	}

	return flags, nil
}

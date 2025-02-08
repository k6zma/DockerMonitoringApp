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
	configFile := flag.String("config_path", "config.json", "Path to json configuration file")
	loggerLevel := flag.String(
		"logger_level",
		"debug",
		"Logger level (debug, info, warn, error, dpanic, panic, fatal)",
	)

	flag.Parse()

	appFlags := &AppFlags{
		ConfigFilePath: *configFile,
		LoggerLevel:    *loggerLevel,
	}

	validate := validator.New()
	if err := validate.Struct(appFlags); err != nil {
		return nil, fmt.Errorf("invalid flags: %w", err)
	}

	return appFlags, nil
}

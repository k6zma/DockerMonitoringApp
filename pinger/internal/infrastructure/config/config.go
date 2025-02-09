package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Ping    *PingConfig    `mapstructure:"ping" validate:"required"`
	Docker  *DockerConfig  `mapstructure:"docker"        validate:"required"`
	Backend *BackendConfig `mapstructure:"backend"       validate:"required"`
}

type BackendConfig struct {
	URL    string `mapstructure:"url"     validate:"required,url"`
	APIKey string `mapstructure:"api_key" validate:"required"`
}

type PingConfig struct {
	PingInterval time.Duration `mapstructure:"ping_interval" validate:"required,gt=4s"`
}

type DockerConfig struct {
	SocketPath string `mapstructure:"socket_path" validate:"required"`
}

func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config read error: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config unmarshal error: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}

	return &cfg, nil
}

package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Server           *ServerConfig     `mapstructure:"server"     validate:"required"`
	DB               *DBConfig         `mapstructure:"db"         validate:"required"`
	MigrationsConfig *MigrationsConfig `mapstructure:"migrations" validate:"required"`
	AuthAPI          *AuthAPIConfig    `mapstructure:"auth_api"   validate:"required"`
}

type ServerConfig struct {
	Port uint16 `mapstructure:"port" validate:"required,gt=0"`
}

type DBConfig struct {
	Host         string `mapstructure:"host"          validate:"required"`
	Port         uint16 `mapstructure:"port"          validate:"required,gt=0"`
	User         string `mapstructure:"user"          validate:"required"`
	Password     string `mapstructure:"password"      validate:"required"`
	DataBaseName string `mapstructure:"database_name" validate:"required"`
}

type MigrationsConfig struct {
	Path string `mapstructure:"path" validate:"required,dir"`
	Type string `mapstructure:"type" validate:"required,oneof=apply drop rollback"`
}

type AuthAPIConfig struct {
	APIKey string `mapstructure:"api_key" validate:"required"`
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

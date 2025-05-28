package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AppMigrate    string `mapstructure:"APP_MIGRATE"`
	AppCasbin     string `mapstructure:"APP_CASBIN"`
	AppPolicy     string `mapstructure:"APP_POLICY"`
	AppHost       string `mapstructure:"APP_HOST"`
	AppPort       string `mapstructure:"APP_PORT"`
	AppPostgres   string `mapstructure:"APP_POSTGRES"`
	AuthMigrate   string `mapstructure:"AUTH_MIGRATE"`
	AuthHost      string `mapstructure:"AUTH_HOST"`
	AuthPort      string `mapstructure:"AUTH_PORT"`
	AuthRedis     string `mapstructure:"AUTH_REDIS"`
	AuthPostgres  string `mapstructure:"AUTH_POSTGRES"`
	AuthSecretKey string `mapstructure:"AUTH_SECRET_KEY"`
}

func New() (*Config, error) {
	viper.AddConfigPath("configs/env")
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	conf := &Config{}
	if err := viper.Unmarshal(conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return conf, nil
}

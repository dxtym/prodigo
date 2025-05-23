package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AuthHost  string `mapstructure:"AUTH_HOST"`
	AuthPort  string `mapstructure:"AUTH_PORT"`
	Redis     string `mapstructure:"REDIS"`
	Postgres  string `mapstructure:"POSTGRES"`
	SecretKey string `mapstructure:"SECRET_KEY"`
}

func New() (*Config, error) {
	viper.AddConfigPath("configs")
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

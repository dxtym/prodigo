package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AuthConfig struct {
	Host      string `mapstructure:"HOST"`
	Port      string `mapstructure:"PORT"`
	Postgres  string `mapstructure:"POSTGRES_DSN"`
	RedisAddr string `mapstructure:"REDIS_ADDR"`
	RedisPass string `mapstructure:"REDIS_PASS"`
	Secret    string `mapstructure:"SECRET_KEY"`
}

func LoadAuthConfig() (*AuthConfig, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("auth")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	conf := &AuthConfig{}
	if err := viper.Unmarshal(conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return conf, nil
}

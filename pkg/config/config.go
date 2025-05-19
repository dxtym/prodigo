package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Host        string `mapstructure:"AUTH_HOST"`
	Port        string `mapstructure:"AUTH_PORT"`
	RedisAddr   string `mapstructure:"REDIS_ADDR"`
	RedisPass   string `mapstructure:"REDIS_PASS"`
	PostgresDSN string `mapstructure:"POSTGRES_DSN"`
	MigrateURL  string `mapstructure:"MIGRATE_URL"`
	Secret      string `mapstructure:"SECRET_KEY"`
}

func LoadAuthConfig() (*Config, error) {
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

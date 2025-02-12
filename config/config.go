package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	GRPCPort    string
	PostgresDSN string
	LogLevel    string
}

// LoadConfig загружает настройки из файла или переменных окружения.
func LoadConfig() (*Config, error) {
	viper.SetConfigFile("config.env") // Указываем путь к конфигурационному файлу
	viper.AutomaticEnv()              // Поддержка переменных окружения

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{
		GRPCPort:    viper.GetString("GRPC_PORT"),
		PostgresDSN: viper.GetString("POSTGRES_DSN"),
		LogLevel:    viper.GetString("LOG_LEVEL"),
	}

	return cfg, nil
}

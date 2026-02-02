package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	Host           string `env:"RUN_ADDRESS"`
	DataBaseDsn    string `env:"DATABASE_URI"`
	SecretKey      string `env:"JWT_SECRET_KEY"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func NewConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}

	flag.StringVar(&cfg.Host, "a", "localhost:8080", "адрес HTTP-сервера")
	flag.StringVar(&cfg.DataBaseDsn, "d", "", "Строка подключения к базе данных")
	flag.StringVar(&cfg.AccrualAddress, "r", "", "адрес системы расчёта начислений")
	flag.StringVar(&cfg.SecretKey, "s", "", "Секретный ключ для JWT")

	flag.Parse()

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

package setup

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/evgenyshipko/golang-metrics-collector/internal/logger"
	"github.com/joho/godotenv"
)

type ServerStartupValues struct {
	Host string `env:"ADDRESS"`
}

func GetStartupValues() ServerStartupValues {
	err := godotenv.Load()
	if err != nil {
		logger.Error(".env файл не найден, используются переменные cистемы")
	}

	host := flag.String("a", "localhost:8080", "input host with port")
	flag.Parse()

	var cfg ServerStartupValues
	err = env.Parse(&cfg)
	if err != nil {
		logger.Error("GetStartupValues: не удалось распарсить env-переменные", err)
	}

	if cfg.Host == "" {
		cfg.Host = *host
	}

	logger.Info(fmt.Sprintf("Параметры запуска: %+v\n", cfg))

	return cfg
}

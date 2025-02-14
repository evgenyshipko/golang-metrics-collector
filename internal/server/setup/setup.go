package setup

import (
	"flag"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"os"
)

type ServerStartupValues struct {
	Host string `env:"ADDRESS"`
}

func GetStartupValues() ServerStartupValues {
	flagHost := flag.String("a", "localhost:8080", "input host with port")
	flag.Parse()

	var cfg ServerStartupValues
	envHost, exists := os.LookupEnv("ADDRESS")

	if exists {
		cfg.Host = envHost
	} else {
		cfg.Host = *flagHost
	}

	logger.Instance.Infow("Параметры запуска:", cfg)

	return cfg
}

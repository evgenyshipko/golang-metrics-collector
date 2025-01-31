package setup

import (
	"flag"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/logger"
	"os"
	"strconv"
)

type AgentStartupValues struct {
	Host           string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

// TODO: функции стартапа в клиенте и сервере похожи. Вот тут бы заюзать мамгию DI, но пока сложно об этом думать
func GetStartupValues() AgentStartupValues {
	flagHost := flag.String("a", "localhost:8080", "metric server host")

	flagReportInterval := flag.Int("r", 10, "interval between report metrics")

	flagPollInterval := flag.Int("p", 2, "interval between polling metrics")

	flag.Parse()

	var cfg AgentStartupValues

	envHost, exists := os.LookupEnv("ADDRESS")
	if exists {
		cfg.Host = envHost
	} else {
		cfg.Host = *flagHost
	}

	envPollInterval, exists := os.LookupEnv("POLL_INTERVAL")
	if exists {
		val, err := strconv.Atoi(envPollInterval)
		if err != nil {
			logger.Error("Ошибка конвертации POLL_INTERVAL")
		}
		cfg.PollInterval = val
	} else {
		cfg.PollInterval = *flagPollInterval
	}

	envReportInterval, exists := os.LookupEnv("REPORT_INTERVAL")
	if exists {
		val, err := strconv.Atoi(envReportInterval)
		if err != nil {
			logger.Error("Ошибка конвертации REPORT_INTERVAL")
		}
		cfg.ReportInterval = val
	} else {
		cfg.ReportInterval = *flagReportInterval
	}

	logger.Info(fmt.Sprintf("Параметры запуска: %+v\n", cfg))

	return cfg
}

package setup

import (
	"errors"
	"flag"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"os"
	"strconv"
	"time"
)

type AgentStartupValues struct {
	Host           string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

func GetStartupValues() (AgentStartupValues, error) {
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

	pollInterval, err := getInterval("POLL_INTERVAL", flagPollInterval)
	if err != nil {
		return AgentStartupValues{}, fmt.Errorf("%w", err)
	}
	cfg.PollInterval = pollInterval

	reportInterval, err := getInterval("REPORT_INTERVAL", flagReportInterval)
	if err != nil {
		return AgentStartupValues{}, fmt.Errorf("%w", err)
	}
	cfg.ReportInterval = reportInterval

	logger.Instance.Infow("Параметры запуска:", cfg)

	return cfg, nil
}

func getInterval(envName string, flagVal *int) (time.Duration, error) {
	envInterval, exists := os.LookupEnv(envName)
	intInterval := 0
	if exists {
		val, err := strconv.Atoi(envInterval)
		if err != nil {
			logger.Instance.Warnw(fmt.Sprintf("ошибка конвертации енва %s, будем драть из флагов", envName))
		}
		intInterval = val
	} else {
		intInterval = *flagVal
	}
	err := validateInterval(intInterval)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	return intToSeconds(intInterval), nil
}

func intToSeconds(num int) time.Duration {
	return time.Duration(num) * time.Second
}

func validateInterval(num int) error {
	if num <= 0 {
		return errors.New("интервал должен быть положительным и больше нуля")
	}
	return nil
}

package setup

import (
	"flag"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/setup"
	"time"
)

type AgentStartupValues struct {
	Host               string          `env:"ADDRESS"`
	ReportInterval     time.Duration   `env:"REPORT_INTERVAL"`
	PollInterval       time.Duration   `env:"POLL_INTERVAL"`
	RetryIntervals     []time.Duration `env:"RETRY_INTERVALS"`
	RequestWaitTimeout time.Duration   `env:"REQUEST_WAIT_TIMEOUT"`
	HashKey            string          `env:"KEY"`
}

const (
	defaultReportIntervalSeconds = 10
	defaultPollIntervalSeconds   = 2
	defaultRequestWaitTimeout    = 10
)

func GetStartupValues(args []string) (AgentStartupValues, error) {
	flagSet := flag.NewFlagSet("config", flag.ContinueOnError)

	flagHost := flagSet.String("a", "localhost:8080", "metric server host")

	flagReportInterval := flagSet.Int("r", defaultReportIntervalSeconds, "interval between report metrics")

	flagPollInterval := flagSet.Int("p", defaultPollIntervalSeconds, "interval between polling metrics")

	flagRetryIntervals := flagSet.String("ri", "1s,3s,5s", "intervals between retries")

	flagRequestWaitTimeout := flagSet.Int("w", defaultRequestWaitTimeout, "http-request wait timeout")

	flagHashKey := flagSet.String("k", "", "secret used to hash metrics")

	// Парсим переданные аргументы
	if err := flagSet.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return AgentStartupValues{}, err
		}
	}

	var cfg AgentStartupValues

	cfg.HashKey = setup.GetStringVariable("KEY", flagHashKey)

	cfg.Host = setup.GetStringVariable("ADDRESS", flagHost)

	pollInterval, err := setup.GetInterval("POLL_INTERVAL", flagPollInterval)
	if err != nil {
		return AgentStartupValues{}, fmt.Errorf("%w", err)
	}
	cfg.PollInterval = pollInterval

	reportInterval, err := setup.GetInterval("REPORT_INTERVAL", flagReportInterval)
	if err != nil {
		return AgentStartupValues{}, fmt.Errorf("%w", err)
	}
	cfg.ReportInterval = reportInterval

	retries, err := setup.GetIntervals("RETRY_INTERVALS", flagRetryIntervals)
	if err != nil {
		return AgentStartupValues{}, err
	}
	cfg.RetryIntervals = retries

	requestWaitTimeout, err := setup.GetInterval("REQUEST_WAIT_TIMEOUT", flagRequestWaitTimeout)
	if err != nil {
		return AgentStartupValues{}, err
	}

	cfg.RequestWaitTimeout = requestWaitTimeout

	logger.Instance.Infow("GetStartupValues", "Параметры запуска:", cfg)

	return cfg, nil
}

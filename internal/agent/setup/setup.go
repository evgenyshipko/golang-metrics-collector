package setup

import (
	"flag"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/setup"
	"time"
)

type AgentStartupValues struct {
	Host           string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

const (
	defaultReportIntervalSeconds = 10
	defaultPollIntervalSeconds   = 2
)

func GetStartupValues() (AgentStartupValues, error) {
	flagHost := flag.String("a", "localhost:8080", "metric server host")

	flagReportInterval := flag.Int("r", defaultReportIntervalSeconds, "interval between report metrics")

	flagPollInterval := flag.Int("p", defaultPollIntervalSeconds, "interval between polling metrics")

	flag.Parse()

	var cfg AgentStartupValues

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

	logger.Instance.Infow("GetStartupValues", "Параметры запуска:", cfg)

	return cfg, nil
}

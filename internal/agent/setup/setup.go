package setup

import (
	"flag"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/files"
	"time"

	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/setup"
)

type AgentStartupValues struct {
	Host                string          `env:"ADDRESS" json:"address"`
	ReportInterval      time.Duration   `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval        time.Duration   `env:"POLL_INTERVAL" json:"poll_interval"`
	RetryIntervals      []time.Duration `env:"RETRY_INTERVALS"`
	RequestWaitTimeout  time.Duration   `env:"REQUEST_WAIT_TIMEOUT"`
	HashKey             string          `env:"KEY"`
	RateLimit           int             `env:"RATE_LIMIT"`
	CryptoPublicKeyPath string          `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigFilePath      string          `env:"CONFIG"`
	Protocol            string          `env:"PROTOCOL"`
}

const (
	defaultReportIntervalSeconds = 5
	defaultPollIntervalSeconds   = 2
	defaultRequestWaitTimeout    = 10
	defaultHostAddress           = "localhost:8080"
	defaultRetryIntervals        = "1s,3s,5s"
	defaultHashKey               = ""
	defaultRateLimit             = 3
	defaultCryptoPublicKeyPath   = ""
	defaultConfigPath            = ""
	defaultProtocol              = "grpc"
)

func GetStartupValues(args []string) (AgentStartupValues, error) {
	flagSet := flag.NewFlagSet("config", flag.ContinueOnError)

	flagHost := flagSet.String("a", defaultHostAddress, "metric server host")

	flagReportInterval := flagSet.Int("r", defaultReportIntervalSeconds, "interval between report metrics")

	flagPollInterval := flagSet.Int("p", defaultPollIntervalSeconds, "interval between polling metrics")

	flagRetryIntervals := flagSet.String("ri", defaultRetryIntervals, "intervals between retries")

	flagRequestWaitTimeout := flagSet.Int("w", defaultRequestWaitTimeout, "http-request wait timeout")

	flagHashKey := flagSet.String("k", defaultHashKey, "secret used to hash metrics")

	flagRateLimit := flagSet.Int("l", defaultRateLimit, "count of http-sender workers")

	flagCryptoPublicKeyPath := flagSet.String("crypto-key", defaultCryptoPublicKeyPath, "path to public key to encrypt metrics")

	flagConfigFilePath := flagSet.String("c", defaultConfigPath, "path to config file")

	flagProtocol := flagSet.String("pr", defaultProtocol, "http or grpc protocol")

	// Парсим переданные аргументы
	if err := flagSet.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return AgentStartupValues{}, err
		}
	}

	var cfg AgentStartupValues

	cfg.Protocol = setup.GetStringVariable("PROTOCOL", flagProtocol)

	cfg.HashKey = setup.GetStringVariable("KEY", flagHashKey)

	cfg.Host = setup.GetStringVariable("ADDRESS", flagHost)

	cfg.CryptoPublicKeyPath = setup.GetStringVariable("CRYPTO_KEY", flagCryptoPublicKeyPath)

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

	rateLimit, err := setup.GetIntVariable("RATE_LIMIT", flagRateLimit)
	if err != nil {
		return AgentStartupValues{}, err
	}

	cfg.RateLimit = rateLimit

	logger.Instance.Infow("GetStartupValues", "Параметры запуска:", cfg)

	// читаем конфиг из json-файла, если он есть
	cfg.ConfigFilePath = setup.GetStringVariable("CONFIG", flagConfigFilePath)
	if cfg.ConfigFilePath != defaultConfigPath {
		configData, err := files.ReadFromFile[AgentStartupValues](cfg.ConfigFilePath)
		if err != nil {
			logger.Instance.Warnw("ReadFromFile", "err", err)
			return cfg, fmt.Errorf("%w", err)
		}
		logger.Instance.Infow("Переменные загружены из json-файла конфига", cfg.ConfigFilePath, configData)
		if configData.Host != defaultHostAddress && cfg.Host == defaultHostAddress {
			cfg.Host = configData.Host
		}
		if configData.ReportInterval != defaultReportIntervalSeconds && cfg.ReportInterval == defaultReportIntervalSeconds {
			cfg.ReportInterval = configData.ReportInterval
		}
		if configData.PollInterval != defaultPollIntervalSeconds && cfg.PollInterval == defaultPollIntervalSeconds {
			cfg.PollInterval = configData.PollInterval
		}
		if configData.CryptoPublicKeyPath != defaultCryptoPublicKeyPath && cfg.CryptoPublicKeyPath == defaultCryptoPublicKeyPath {
			cfg.CryptoPublicKeyPath = configData.CryptoPublicKeyPath
		}
	}

	return cfg, nil
}

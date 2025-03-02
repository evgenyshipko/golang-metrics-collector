package setup

import (
	"flag"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/setup"
	"os"
	"path/filepath"
	"time"
)

type ServerStartupValues struct {
	Host               string          `env:"ADDRESS"`
	StoreInterval      time.Duration   `env:"STORE_INTERVAL"`
	FileStoragePath    string          `env:"FILE_STORAGE_PATH"`
	Restore            bool            `env:"RESTORE"`
	DatabaseDSN        string          `env:"DATABASE_DSN"`
	RetryIntervals     []time.Duration `env:"RETRY_INTERVALS"`
	RequestWaitTimeout time.Duration   `env:"REQUEST_WAIT_TIMEOUT"`
	AutoMigrations     bool            `env:"AUTO_MIGRATIONS"`
}

func GetProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		logger.Instance.Fatal(err.Error())
	}
	return dir
}

const (
	defaultStoreIntervalSeconds = 300
	defaultRequestWaitTimeout   = 10
)

func GetStartupValues(args []string) (ServerStartupValues, error) {

	rootDir := GetProjectRoot()
	defaultFilePath := filepath.Join(rootDir, "temp.json")

	flagSet := flag.NewFlagSet("config", flag.ContinueOnError)

	flagHost := flagSet.String("a", "localhost:8080", "input host with port")

	flagStoreInterval := flagSet.Int("i", defaultStoreIntervalSeconds, "interval between saving metrics to file")

	flagFileStoragePath := flagSet.String("f", defaultFilePath, "temp file to store metrics")

	flagRestore := flagSet.Bool("r", true, "restore saved metrics from file or not")

	// postgres://metrics:metrics@localhost:5433/metrics?sslmode=disable&connect_timeout=5
	flagDatabaseDsn := flagSet.String("d", "", "database dsn")

	flagRetryIntervals := flagSet.String("ri", "1s,3s,5s", "intervals between retries")

	flagRequestWaitTimeout := flagSet.Int("w", defaultRequestWaitTimeout, "http-request wait timeout")

	flagAutoMigrations := flagSet.Bool("m", true, "run migrations on server startup")

	// Парсим переданные аргументы
	if err := flagSet.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return ServerStartupValues{}, err
		}
	}

	var cfg ServerStartupValues

	cfg.DatabaseDSN = setup.GetStringVariable("DATABASE_DSN", flagDatabaseDsn)

	cfg.Host = setup.GetStringVariable("ADDRESS", flagHost)

	cfg.FileStoragePath = setup.GetStringVariable("FILE_STORAGE_PATH", flagFileStoragePath)

	cfg.Restore = setup.GetBoolVariable("RESTORE", flagRestore)

	storeInterval, err := setup.GetInterval("STORE_INTERVAL", flagStoreInterval, false)
	if err != nil {
		return ServerStartupValues{}, err
	}
	cfg.StoreInterval = storeInterval

	retries, err := setup.GetIntervals("RETRY_INTERVALS", flagRetryIntervals)
	if err != nil {
		return ServerStartupValues{}, err
	}
	cfg.RetryIntervals = retries

	requestWaitTimeout, err := setup.GetInterval("REQUEST_WAIT_TIMEOUT", flagRequestWaitTimeout)
	if err != nil {
		return ServerStartupValues{}, err
	}

	cfg.RequestWaitTimeout = requestWaitTimeout

	cfg.AutoMigrations = setup.GetBoolVariable("AUTO_MIGRATIONS", flagAutoMigrations)

	logger.Instance.Infow("GetStartupValues", "Параметры запуска:", cfg)

	return cfg, nil
}

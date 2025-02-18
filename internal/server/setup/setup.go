package setup

import (
	"flag"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/setup"
	"time"
)

type ServerStartupValues struct {
	Host            string        `env:"ADDRESS"`
	StoreInterval   time.Duration `env:"STORE_INTERVAL"`
	FileStoragePath string        `env:"FILE_STORAGE_PATH"`
	Restore         bool          `env:"RESTORE"`
}

func GetStartupValues() (ServerStartupValues, error) {
	flagHost := flag.String("a", "localhost:8080", "input host with port")

	flagStoreInterval := flag.Int("i", 300, "interval between saving metrics to file")

	flagFileStoragePath := flag.String("f", "./temp.json", "temp file to store metrics")

	flagRestore := flag.Bool("r", true, "restore saved metrics from file or not")

	flag.Parse()

	var cfg ServerStartupValues

	cfg.Host = setup.GetStringVariable("ADDRESS", flagHost)

	cfg.FileStoragePath = setup.GetStringVariable("FILE_STORAGE_PATH", flagFileStoragePath)

	cfg.Restore = setup.GetBoolVariable("RESTORE", flagRestore)

	storeInterval, err := setup.GetInterval("STORE_INTERVAL", flagStoreInterval)
	if err != nil {
		return ServerStartupValues{}, fmt.Errorf("%w", err)
	}
	cfg.StoreInterval = storeInterval

	logger.Instance.Infow("GetStartupValues", "Параметры запуска:", cfg)

	return cfg, nil
}

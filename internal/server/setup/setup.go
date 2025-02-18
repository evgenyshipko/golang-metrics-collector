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

func GetStartupValues(args []string) (ServerStartupValues, error) {

	flagSet := flag.NewFlagSet("config", flag.ContinueOnError)

	flagHost := flagSet.String("a", "localhost:8080", "input host with port")

	flagStoreInterval := flagSet.Int("i", 300, "interval between saving metrics to file")

	flagFileStoragePath := flagSet.String("f", "./temp.json", "temp file to store metrics")

	flagRestore := flagSet.Bool("r", true, "restore saved metrics from file or not")

	// Парсим переданные аргументы
	if err := flagSet.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return ServerStartupValues{}, err
		}
	}

	var cfg ServerStartupValues

	cfg.Host = setup.GetStringVariable("ADDRESS", flagHost)

	cfg.FileStoragePath = setup.GetStringVariable("FILE_STORAGE_PATH", flagFileStoragePath)

	cfg.Restore = setup.GetBoolVariable("RESTORE", flagRestore)

	storeInterval, err := setup.GetInterval("STORE_INTERVAL", flagStoreInterval, false)
	if err != nil {
		return ServerStartupValues{}, fmt.Errorf("%w", err)
	}
	cfg.StoreInterval = storeInterval

	logger.Instance.Infow("GetStartupValues", "Параметры запуска:", cfg)

	return cfg, nil
}

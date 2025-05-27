// setup - пакет с настройками запуска приложения. Приложение может брать переменные из флагов командной строки или из .env-файла.
package setup

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/setup"
)

type ServerStartupValues struct {
	Host               string          `env:"ADDRESS"`              // Host определяет адрес и порт, на котором сервер будет слушать входящие соединения (флаг -a).
	StoreInterval      time.Duration   `env:"STORE_INTERVAL"`       // Интервал времени между сохранениями метрик в локальный файл (флаг -i).
	FileStoragePath    string          `env:"FILE_STORAGE_PATH"`    // Путь к файлу, в котором сохраняются метрики (флаг -f).
	Restore            bool            `env:"RESTORE"`              // Восстанавливать метрики из файла при запуске или нет (флаг -r).
	DatabaseDSN        string          `env:"DATABASE_DSN"`         // Строка с данными доступа к базе PostgreSQL (флаг -d).
	RetryIntervals     []time.Duration `env:"RETRY_INTERVALS"`      // Интервалы между попытками записи в базу (флаг -ri).
	RequestWaitTimeout time.Duration   `env:"REQUEST_WAIT_TIMEOUT"` // Таймаут ожидания ответа хендлеров (флаг -w).
	AutoMigrations     bool            `env:"AUTO_MIGRATIONS"`      // Запускать миграции при запуске приложения или нет (флаг -m).
	HashKey            string          `env:"KEY"`                  // Секретный хеш для авторизации (флаг -k).
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

// GetStartupValues берет переменные из флагов либо из переменных окружения. Если нет ни того, ни другого - то берет дефолтные значения.
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

	flagHashKey := flagSet.String("k", "", "secret used to hash metrics")

	// Парсим переданные аргументы
	if err := flagSet.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return ServerStartupValues{}, err
		}
	}

	var cfg ServerStartupValues

	cfg.HashKey = setup.GetStringVariable("KEY", flagHashKey)

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

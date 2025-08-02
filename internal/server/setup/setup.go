// setup - пакет с настройками запуска приложения. Приложение может брать переменные из флагов командной строки или из .env-файла.
package setup

import (
	"flag"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/files"
	"os"
	"path/filepath"
	"time"

	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/setup"
)

type ServerStartupValues struct {
	Host                 string          `env:"ADDRESS" json:"address"`               // Host определяет адрес и порт, на котором сервер будет слушать входящие соединения (флаг -a).
	StoreInterval        time.Duration   `env:"STORE_INTERVAL" json:"store_interval"` // Интервал времени между сохранениями метрик в локальный файл (флаг -i).
	FileStoragePath      string          `env:"FILE_STORAGE_PATH" json:"store_file"`  // Путь к файлу, в котором сохраняются метрики (флаг -f).
	Restore              bool            `env:"RESTORE" json:"restore"`               // Восстанавливать метрики из файла при запуске или нет (флаг -r).
	DatabaseDSN          string          `env:"DATABASE_DSN" json:"database_dsn"`     // Строка с данными доступа к базе PostgreSQL (флаг -d).
	RetryIntervals       []time.Duration `env:"RETRY_INTERVALS"`                      // Интервалы между попытками записи в базу (флаг -ri).
	RequestWaitTimeout   time.Duration   `env:"REQUEST_WAIT_TIMEOUT"`                 // Таймаут ожидания ответа хендлеров (флаг -w).
	AutoMigrations       bool            `env:"AUTO_MIGRATIONS"`                      // Запускать миграции при запуске приложения или нет (флаг -m).
	HashKey              string          `env:"KEY"`                                  // Секретный хеш для авторизации (флаг -k).
	CryptoPrivateKeyPath string          `env:"CRYPTO_KEY" json:"crypto_key"`         // Путь к секретному приватному ключу для расшифровки сообщений, подписанных публичным ключом шифрования (флаг -crypto-key)
	ConfigFilePath       string          `env:"CONFIG"`
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
	defaultHostAddress          = "localhost:8080"
	defaultRestoreFlag          = true
	defaultDatabaseRSN          = ""
	defaultRetryIntervals       = "1s,3s,5s"
	defaultAutoMigrationsFlag   = true
	defaultHashKey              = ""
	defaultCryptoPrivateKeyPath = ""
	defaultConfigPath           = ""
)

// GetStartupValues берет переменные из флагов либо из переменных окружения. Если нет ни того, ни другого - то берет дефолтные значения.
func GetStartupValues(args []string) (ServerStartupValues, error) {

	rootDir := GetProjectRoot()
	defaultFileStoragePath := filepath.Join(rootDir, "temp.json")

	flagSet := flag.NewFlagSet("config", flag.ContinueOnError)

	flagHost := flagSet.String("a", defaultHostAddress, "input host with port")

	flagStoreInterval := flagSet.Int("i", defaultStoreIntervalSeconds, "interval between saving metrics to file")

	flagFileStoragePath := flagSet.String("f", defaultFileStoragePath, "temp file to store metrics")

	flagRestore := flagSet.Bool("r", defaultRestoreFlag, "restore saved metrics from file or not")

	// postgres://metrics:metrics@localhost:5433/metrics?sslmode=disable&connect_timeout=5
	flagDatabaseDsn := flagSet.String("d", defaultDatabaseRSN, "database dsn")

	flagRetryIntervals := flagSet.String("ri", defaultRetryIntervals, "intervals between retries")

	flagRequestWaitTimeout := flagSet.Int("w", defaultRequestWaitTimeout, "http-request wait timeout")

	flagAutoMigrations := flagSet.Bool("m", defaultAutoMigrationsFlag, "run migrations on server startup")

	flagHashKey := flagSet.String("k", defaultHashKey, "secret used to hash metrics")

	cryptoPrivateKeyPath := flagSet.String("crypto-key", defaultCryptoPrivateKeyPath, "path to public key to encrypt metrics")

	flagConfigFilePath := flagSet.String("c", defaultConfigPath, "path to config file")

	// Парсим переданные аргументы
	if err := flagSet.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return ServerStartupValues{}, err
		}
	}

	var cfg ServerStartupValues

	cfg.CryptoPrivateKeyPath = setup.GetStringVariable("CRYPTO_KEY", cryptoPrivateKeyPath)

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

	// читаем конфиг из json-файла, если он есть
	cfg.ConfigFilePath = setup.GetStringVariable("CONFIG", flagConfigFilePath)
	if cfg.ConfigFilePath != defaultConfigPath {
		configData, err := files.ReadFromFile[ServerStartupValues](cfg.ConfigFilePath)
		if err != nil {
			logger.Instance.Warnw("ReadFromFile", "err", err)
			return cfg, fmt.Errorf("%w", err)
		}
		logger.Instance.Infow("Переменные загружены из json-файла конфига", cfg.ConfigFilePath, configData)
		if configData.Host != defaultHostAddress && cfg.Host == defaultHostAddress {
			cfg.Host = configData.Host
		}
		if configData.Restore != defaultRestoreFlag && cfg.Restore == defaultRestoreFlag {
			cfg.Restore = configData.Restore
		}
		if configData.StoreInterval != defaultStoreIntervalSeconds && cfg.StoreInterval == defaultStoreIntervalSeconds {
			cfg.StoreInterval = configData.StoreInterval
		}
		if configData.FileStoragePath != defaultFileStoragePath && cfg.FileStoragePath == defaultFileStoragePath {
			cfg.FileStoragePath = configData.FileStoragePath
		}
		if configData.DatabaseDSN != defaultDatabaseRSN && cfg.DatabaseDSN == defaultDatabaseRSN {
			cfg.DatabaseDSN = configData.DatabaseDSN
		}
		if configData.CryptoPrivateKeyPath != defaultCryptoPrivateKeyPath && cfg.CryptoPrivateKeyPath == defaultCryptoPrivateKeyPath {
			cfg.CryptoPrivateKeyPath = configData.CryptoPrivateKeyPath
		}
	}

	logger.Instance.Infow("GetStartupValues", "Параметры запуска:", cfg)

	return cfg, nil
}

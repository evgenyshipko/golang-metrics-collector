package files

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/retry"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
)

func ReadFromFile(fileName string) (*storage.StorageData, error) {
	logger.Instance.Info("Читаем из файла")
	consumer, err := NewConsumer(fileName)
	if err != nil {
		logger.Instance.Warnw("ReadFromFile", "NewConsumer", err)
		return nil, err
	}
	storageData, err := consumer.ReadData()
	if err != nil {
		logger.Instance.Warnw("ReadFromFile", "consumer.ReadData", err)
		return nil, err
	}
	logger.Instance.Infow("ReadFromFile", "Прочитано успешно", *storageData)
	return storageData, nil
}

func ReadFromFileWithRetry(fileName string) (*storage.StorageData, error) {
	return retry.WithRetry(func() (*storage.StorageData, error) {
		return ReadFromFile(fileName)
	})
}

package files

import (
	"time"

	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/retry"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
)

func WriteToFile(fileName string, data *storage.StorageData) error {
	logger.Instance.Info("Пишем метрики в файл")

	producer, err := NewTruncateProducer(fileName)
	if err != nil {
		logger.Instance.Warnw("WriteToFile", "NewTruncateProducer", err)
		return err
	}

	err = producer.WriteData(data)
	if err != nil {
		logger.Instance.Warnw("WriteToFile", "WriteData", err)
		return err
	}
	logger.Instance.Infof("Файл %s записан успешно", fileName)
	return err
}

func WriteToFileWithRetry(fileName string, data *storage.StorageData, retryIntervals []time.Duration) error {
	_, err := retry.WithRetry(func() (string, error) {
		err := WriteToFile(fileName, data)
		return "", err
	}, retryIntervals)
	return err
}

package files

import (
	"time"

	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/retry"
)

func WriteToFile[T any](fileName string, data *T) error {
	logger.Instance.Info("Пишем метрики в файл")

	producer, err := NewTruncateProducer(fileName)
	if err != nil {
		logger.Instance.Warnw("WriteToFile", "NewTruncateProducer", err)
		return err
	}

	err = WriteData(producer, data)
	if err != nil {
		logger.Instance.Warnw("WriteToFile", "WriteData", err)
		return err
	}
	logger.Instance.Infof("Файл %s записан успешно", fileName)
	return err
}

func WriteToFileWithRetry[T any](fileName string, data *T, retryIntervals []time.Duration) error {
	_, err := retry.WithRetry(func() (string, error) {
		err := WriteToFile(fileName, data)
		return "", err
	}, retryIntervals)
	return err
}

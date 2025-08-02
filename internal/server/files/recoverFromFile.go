package files

import (
	"context"
	files2 "github.com/evgenyshipko/golang-metrics-collector/internal/common/files"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/retry"
	"time"

	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
)

func RecoverFromFile(filePath string, store storage.Storage, retryIntervals []time.Duration) {
	fileData, err := ReadFromFileWithRetry(filePath, retryIntervals)
	if err != nil {
		logger.Instance.Warnw("ReadFromFileWithRetry", "Ошибка чтения из файла", err)
		return
	}
	err = store.SetData(context.Background(), *fileData)
	if err != nil {
		logger.Instance.Warnw("ReadFromFileWithRetry", "storage.SetData ошибка записи в хранилище", err)
		return
	}
	logger.Instance.Infof("Хранилище восстановлено из файла %s успешно", filePath)
}

func ReadFromFileWithRetry(fileName string, retryIntervals []time.Duration) (*storage.StorageData, error) {
	return retry.WithRetry(func() (*storage.StorageData, error) {
		return files2.ReadFromFile[storage.StorageData](fileName)
	}, retryIntervals)
}

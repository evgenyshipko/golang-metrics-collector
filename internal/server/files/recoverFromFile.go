package files

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
)

func RecoverFromFile(filePath string, store storage.Storage) {
	fileData, err := ReadFromFileWithRetry(filePath)
	if err != nil {
		logger.Instance.Warnw("ReadFromFileWithRetry", "Ошибка чтения из файла", err)
		return
	}
	err = store.SetData(*fileData)
	if err != nil {
		logger.Instance.Warnw("ReadFromFileWithRetry", "storage.SetData ошибка записи в хранилище", err)
		return
	}
	logger.Instance.Infof("Хранилище восстановлено из файла %s успешно", filePath)
}

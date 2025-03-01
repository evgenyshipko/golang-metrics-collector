package files

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
)

func ReadFromFile(fileName string, storage storage.Storage) {
	logger.Instance.Info("Читаем из файла")
	consumer, err := NewConsumer(fileName)
	if err != nil {
		logger.Instance.Warnw("ReadFromFile", "NewConsumer", err)
		return
	}
	storageData, err := consumer.ReadData()
	if err != nil {
		logger.Instance.Warnw("ReadFromFile", "consumer.ReadData", err)
		return
	}
	err = storage.SetData(*storageData)
	if err != nil {
		logger.Instance.Warnw("ReadFromFile", "storage.SetData ошибка записи в хранилище", err)
	}
	logger.Instance.Infow("ReadFromFile", "Прочитано успешно", *storageData)
}

package files

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
)

func ReadFromFile(fileName string, storage *storage.MemStorage) {
	logger.Instance.Info("Читаем из файла")
	consumer, err := NewConsumer(fileName)
	if err != nil {
		logger.Instance.Warnw("ReadFromFile", "NewConsumer", err)
		return
	}
	memStorageData, err := consumer.ReadData()
	if err != nil {
		logger.Instance.Warnw("ReadFromFile", "consumer.ReadData", err)
		return
	}
	storage.SetData(*memStorageData)
	logger.Instance.Infow("ReadFromFile", "Прочитано успешно", *memStorageData)
}

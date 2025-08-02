package files

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
)

func ReadFromFile[T any](fileName string) (*T, error) {
	logger.Instance.Info("Читаем из файла")
	consumer, err := NewConsumer(fileName)
	if err != nil {
		logger.Instance.Warnw("ReadFromFile", "NewConsumer", err)
		return nil, err
	}
	storageData, err := ReadData[T](consumer)
	if err != nil {
		logger.Instance.Warnw("ReadFromFile", "consumer.ReadData", err)
		return nil, err
	}
	logger.Instance.Infow("ReadFromFile", "Прочитано успешно", *storageData)
	return storageData, nil
}

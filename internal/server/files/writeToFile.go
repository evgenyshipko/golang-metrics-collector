package files

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
)

func WriteToFile(fileName string, data *storage.StorageData) {
	logger.Instance.Info("Пишем метрики в файл")

	producer, err := NewTruncateProducer(fileName)
	if err != nil {
		logger.Instance.Warnw("StoreMetricHandler", "NewTruncateProducer", err)
	}

	err = producer.WriteData(data)
	if err != nil {
		logger.Instance.Warnw("StoreMetricHandler", "WriteData", err)
	}
	logger.Instance.Info("Файл записан успешно")
}

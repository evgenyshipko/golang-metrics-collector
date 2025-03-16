package tasks

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/storage"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/types"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
)

func MetricsConsumer(store *storage.MetricStorage, dataChan <-chan types.ChanData) {
	for {
		data := <-dataChan
		if data.Err != nil {
			logger.Instance.Warnw("Обработка ошибки", "error", data.Err)
			continue
		}
		logger.Instance.Infow("Принято сообщение в канале", "Name", data.Data.Name, "Type", data.Data.Type, "Value", data.Data.Value)
		store.Set(data.Data)
	}
}

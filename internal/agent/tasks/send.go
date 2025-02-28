package tasks

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/converter"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/requests"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/storage"
	c "github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"time"
)

func SendMetricsTask(interval time.Duration, metrics *storage.MetricStorage, host string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {

		var metricDataArr []c.MetricData

		for name, data := range *metrics {
			metricData, err := converter.GenerateMetricData(data.Type, name, data.Value)
			if err != nil {
				logger.Instance.Warnw("SendMetricsTask", "GenerateMetricData", err)
				return
			}
			metricDataArr = append(metricDataArr, metricData)
		}

		err := requests.SendMetricBatch(host, metricDataArr)
		if err != nil {
			logger.Instance.Warnw("SendMetricsTask", "SendMetricBatch", err)
			return
		}
	}
}

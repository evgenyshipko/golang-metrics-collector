package tasks

import (
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/requests"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/storage"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"time"
)

func SendMetricsTask(interval time.Duration, metrics *storage.MetricStorage, host string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		for metricName, metricData := range *metrics {
			err := requests.SendMetric(host, metricData.Type, metricName, metricData.Value)
			if err != nil {
				logger.Instance.Warnw(fmt.Sprintf("SendMetricsTask %s", err))
				continue
			}
		}
	}
}

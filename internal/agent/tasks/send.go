package tasks

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/converter"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/requests"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/storage"
	c "github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"time"
)

func SendMetricsTask(cfg setup.AgentStartupValues, storage *storage.MetricStorage) {
	ticker := time.NewTicker(cfg.ReportInterval)
	defer ticker.Stop()

	requester := requests.NewRequester(cfg)

	for range ticker.C {

		var metricDataArr []c.MetricData

		for name, data := range storage.Data {
			metricData, err := converter.GenerateMetricData(data.Type, name, data.Value)
			if err != nil {
				logger.Instance.Warnw("SendMetricsTask", "GenerateMetricData", err)
				return
			}
			metricDataArr = append(metricDataArr, metricData)
		}

		err := requester.SendMetricBatch(cfg.Host, metricDataArr)
		if err != nil {
			logger.Instance.Warnw("SendMetricsTask", "SendMetricBatch", err)
			return
		}
	}
}

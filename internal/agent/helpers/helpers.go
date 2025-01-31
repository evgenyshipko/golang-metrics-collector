package helpers

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/consts"
)

// TODO: вместо этой штуки лучше хранить рядом с именем метрики ее тип. Было лень сразу сделано нормально, надо переделать.
func GetMetricType(metricName string) consts.Metric {
	if metricName == "PollCount" {
		return consts.COUNTER
	}
	return consts.GAUGE
}

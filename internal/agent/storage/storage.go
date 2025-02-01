package storage

import "github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"

type MetricData struct {
	Value interface{}
	Type  consts.Metric
}

type MetricStorage map[string]MetricData

func NewMetricStorage() MetricStorage {
	return make(MetricStorage)
}

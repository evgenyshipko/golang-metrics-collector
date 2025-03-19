package types

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
)

type MetricValue struct {
	Type  consts.Metric
	Value interface{}
	Name  string
}

type MetricMessage struct {
	Data MetricValue
	Err  error
}

type MetricData struct {
	Value interface{}
	Type  consts.Metric
}

type MetricDataMap map[string]MetricData

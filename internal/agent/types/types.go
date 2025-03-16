package types

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
)

type Data struct {
	Type  consts.Metric
	Value interface{}
	Name  string
}

type ChanData struct {
	Data Data
	Err  error
}

type MetricData struct {
	Value interface{}
	Type  consts.Metric
}

type MetricDataMap map[string]MetricData

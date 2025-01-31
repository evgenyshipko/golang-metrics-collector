package consts

type Metric string

const (
	GAUGE   Metric = "gauge"
	COUNTER Metric = "counter"
)

type Gauge float64
type Counter int64

const (
	METRIC_TYPE  string = "metricType"
	METRIC_VALUE string = "metricValue"
	METRIC_NAME  string = "metricName"
)

package consts

type Metric string

const (
	GAUGE   Metric = "gauge"
	COUNTER Metric = "counter"
)

type Gauge float64
type Counter int64

const (
	MetricType  string = "metricType"
	MetricValue string = "metricValue"
	MetricName  string = "metricName"
)

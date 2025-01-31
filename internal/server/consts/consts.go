package consts

type Metric string

const (
	GAUGE   Metric = "gauge"
	COUNTER Metric = "counter"
)

type Gauge float64
type Counter int64

type UrlParam string

const (
	MetricType  UrlParam = "metricType"
	MetricValue UrlParam = "metricValue"
	MetricName  UrlParam = "metricName"
)

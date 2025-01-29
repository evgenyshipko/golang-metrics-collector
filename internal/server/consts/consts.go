package consts

type Metric string

const (
	GAUGE   Metric = "gauge"
	COUNTER Metric = "counter"
)

type Gauge float64
type Counter int64

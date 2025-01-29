package consts

type Metrics string

const (
	GAUGE   Metrics = "gauge"
	COUNTER Metrics = "counter"
)

type Gauge float64
type Counter int64

package consts

type Metric string

const (
	GAUGE   Metric = "gauge"
	COUNTER Metric = "counter"
)

type MetricData struct {
	ID    string   `json:"id"`              // имя метрики
	MType Metric   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewMetricData(metricType Metric, name string, value Values) (MetricData, error) {
	return MetricData{
		ID:    name,
		MType: metricType,
		Value: value.Gauge,
		Delta: value.Counter,
	}, nil
}

type URLParam string

const (
	MetricType  URLParam = "metricType"
	MetricValue URLParam = "metricValue"
	MetricName  URLParam = "metricName"
)

type Values struct {
	Counter *int64   `json:"counter,omitempty"`
	Gauge   *float64 `json:"gauge,omitempty"`
}

const HashSha256Header = "HashSHA256"

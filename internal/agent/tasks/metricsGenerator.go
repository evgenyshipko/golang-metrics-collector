package tasks

import (
	"time"

	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/collector"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/types"
)

func MetricsGenerator(interval time.Duration, inputCh chan<- types.MetricMessage) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	collectMetricsMemoized := collector.Create()

	for range ticker.C {
		metrics := *collectMetricsMemoized()

		for name, value := range metrics {
			inputCh <- types.MetricMessage{
				Data: types.MetricValue{
					Type:  value.Type,
					Name:  name,
					Value: value.Value,
				},
				Err: nil,
			}
		}
	}
}

package tasks

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/collector"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/types"
	"time"
)

func MetricsGenerator(interval time.Duration, inputCh chan<- types.ChanData) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	collectMetricsMemoized := collector.Create()

	for range ticker.C {
		metrics := *collectMetricsMemoized()

		for name, value := range metrics {
			inputCh <- types.ChanData{
				Data: types.Data{
					Type:  value.Type,
					Name:  name,
					Value: value.Value,
				},
				Err: nil,
			}
		}
	}
}

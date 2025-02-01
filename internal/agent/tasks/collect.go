package tasks

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/collector"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/storage"
	"time"
)

func CollectMetricsTask(interval time.Duration, metrics *storage.MetricStorage) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	collectMetricsMemoized := collector.Create(metrics)

	for range ticker.C {
		collectMetricsMemoized()
	}
}

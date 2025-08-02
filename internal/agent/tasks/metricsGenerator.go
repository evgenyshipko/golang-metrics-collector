package tasks

import (
	"context"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"time"

	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/collector"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/types"
)

func MetricsGenerator(ctx context.Context, interval time.Duration, inputCh chan<- types.MetricMessage) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	collectMetricsMemoized := collector.Create()

	for {
		select {
		case <-ctx.Done(): // Получен сигнал завершения
			logger.Instance.Debug("MetricsGenerator between ticker <-ctx.Done()")
			return
		case <-ticker.C:
			metrics := *collectMetricsMemoized()

			for name, value := range metrics {
				// Проверяем контекст перед отправкой каждой метрики
				select {
				case <-ctx.Done():
					logger.Instance.Debug("MetricsGenerator before sending <-ctx.Done()")
					return
				case inputCh <- types.MetricMessage{
					Data: types.MetricValue{
						Type:  value.Type,
						Name:  name,
						Value: value.Value,
					},
					Err: nil,
				}:
					// Метрика успешно отправлена
				}
			}
		}
	}
}

package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/types"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

/*
Ключевые особенности реализации:
1 Контроль контекста в нескольких точках:
- Перед началом сбора метрик
-Перед каждой отправкой в канал
-В начале каждой итерации цикла

2 Неблокирующая отправка метрик:
-Все операции отправки в канал защищены select с проверкой ctx.Done()
-Гарантирует быструю реакцию на сигнал завершения
*/

func AdditionalMetricsGenerator(ctx context.Context, interval time.Duration, inputCh chan<- types.MetricMessage) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done(): // Получен сигнал завершения
			return
		case <-ticker.C:
			if ctx.Err() != nil {
				return
			}

			memory, err := mem.VirtualMemory()
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case inputCh <- types.MetricMessage{
					Data: types.MetricValue{},
					Err:  err,
				}:
				}
				continue
			}

			memMetrics := []types.MetricMessage{
				{
					Data: types.MetricValue{
						Type:  consts.GAUGE,
						Value: memory.Total,
						Name:  "TotalMemory",
					},
					Err: nil,
				},
				{
					Data: types.MetricValue{
						Type:  consts.GAUGE,
						Value: memory.Free,
						Name:  "FreeMemory",
					},
					Err: nil,
				},
			}

			for _, metric := range memMetrics {
				select {
				case <-ctx.Done():
					return
				case inputCh <- metric:
				}
			}

			cpuPercent, err := cpu.Percent(interval, true)
			if err != nil {
				logger.Instance.Warnf("Ошибка при получении загрузки CPU: %v\n", err)
				select {
				case <-ctx.Done():
					return
				case inputCh <- types.MetricMessage{
					Data: types.MetricValue{},
					Err:  err,
				}:
				}
				continue
			}

			for index, value := range cpuPercent {
				select {
				case <-ctx.Done():
					return
				case inputCh <- types.MetricMessage{
					Data: types.MetricValue{
						Type:  consts.GAUGE,
						Value: value,
						Name:  fmt.Sprintf("CPUutilization%d", index+1),
					},
					Err: nil,
				}:
				}
			}
		}
	}
}

package tasks

import (
	"fmt"
	"time"

	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/types"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func AdditionalMetricsGenerator(interval time.Duration, inputCh chan<- types.MetricMessage) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		memory, err := mem.VirtualMemory()
		if err != nil {
			inputCh <- types.MetricMessage{
				Data: types.MetricValue{},
				Err:  err,
			}
		} else {
			inputCh <- types.MetricMessage{
				Data: types.MetricValue{
					Type:  consts.GAUGE,
					Value: memory.Total,
					Name:  "TotalMemory",
				},
				Err: nil,
			}
			inputCh <- types.MetricMessage{
				Data: types.MetricValue{
					Type:  consts.GAUGE,
					Value: memory.Free,
					Name:  "FreeMemory",
				},
				Err: nil,
			}
		}

		cpuPercent, err := cpu.Percent(interval, true)
		if err != nil {
			logger.Instance.Warnf("Ошибка при получении загрузки CPU: %v\n", err)
			inputCh <- types.MetricMessage{
				Data: types.MetricValue{},
				Err:  err,
			}
		} else {
			for index, value := range cpuPercent {
				inputCh <- types.MetricMessage{
					Data: types.MetricValue{
						Type:  consts.GAUGE,
						Value: value,
						Name:  fmt.Sprintf("CPUutilization%d", index+1),
					},
					Err: nil,
				}
			}

		}
	}
}

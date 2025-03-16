package tasks

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/types"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"time"
)

func AdditionalMetricsGenerator(interval time.Duration, inputCh chan<- types.ChanData) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {

		memory, err := mem.VirtualMemory()
		if err != nil {
			inputCh <- types.ChanData{
				Data: types.Data{},
				Err:  err,
			}
		} else {
			inputCh <- types.ChanData{
				Data: types.Data{
					Type:  consts.GAUGE,
					Value: memory.Total,
					Name:  "TotalMemory",
				},
				Err: nil,
			}
			inputCh <- types.ChanData{
				Data: types.Data{
					Type:  consts.GAUGE,
					Value: memory.Free,
					Name:  "FreeMemory",
				},
				Err: nil,
			}
		}

		cpuPercent, err := cpu.Percent(time.Minute, false)
		if err != nil {
			logger.Instance.Warnf("Ошибка при получении загрузки CPU: %v\n", err)
			inputCh <- types.ChanData{
				Data: types.Data{},
				Err:  err,
			}
		} else {
			inputCh <- types.ChanData{
				Data: types.Data{
					Type:  consts.GAUGE,
					Value: cpuPercent[0],
					Name:  "CPUutilization1",
				},
				Err: nil,
			}
		}
	}
}

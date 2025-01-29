package tasks

import (
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/helpers"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/requests"
	"math/rand"
	"runtime"
	"time"
)

func CollectMetricsTask(interval time.Duration, metrics *map[string]interface{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	PollCount := 0

	for {
		select {
		case <-ticker.C:
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			fmt.Println("Сбор и сохранение метрик")
			PollCount++
			*metrics = map[string]interface{}{
				"Alloc":         memStats.Alloc,
				"BuckHashSys":   memStats.BuckHashSys,
				"Frees":         memStats.Frees,
				"GCCPUFraction": memStats.GCCPUFraction,
				"GCSys":         memStats.GCSys,
				"HeapAlloc":     memStats.HeapAlloc,
				"HeapIdle":      memStats.HeapIdle,
				"HeapInuse":     memStats.HeapInuse,
				"HeapObjects":   memStats.HeapObjects,
				"HeapReleased":  memStats.HeapReleased,
				"HeapSys":       memStats.HeapSys,
				"LastGC":        memStats.LastGC,
				"Lookups":       memStats.Lookups,
				"MCacheInuse":   memStats.MCacheInuse,
				"MCacheSys":     memStats.MCacheSys,
				"MSpanInuse":    memStats.MSpanInuse,
				"MSpanSys":      memStats.MSpanSys,
				"Mallocs":       memStats.Mallocs,
				"NextGC":        memStats.NextGC,
				"NumForcedGC":   memStats.NumForcedGC,
				"NumGC":         memStats.NumGC,
				"OtherSys":      memStats.OtherSys,
				"PauseTotalNs":  memStats.PauseTotalNs,
				"StackInuse":    memStats.StackInuse,
				"StackSys":      memStats.StackSys,
				"Sys":           memStats.Sys,
				"TotalAlloc":    memStats.TotalAlloc,
				"PollCount":     PollCount,
				"RandomValue":   rand.Float64(),
			}
			fmt.Println("Данные", metrics)
		}
	}
}

func SendMetricsTask(interval time.Duration, metrics *map[string]interface{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for metricName, metricValue := range *metrics {

				value, err := helpers.ConvertMetricValueToString(metricName, metricValue)
				if err != nil {
					fmt.Println("ERROR in ConvertMetricValueToString", err)
					continue
				}

				err = requests.SendMetric(helpers.GetMetricType(metricName), metricName, value)
				if err != nil {
					fmt.Println("ERROR in sendMetricsTask", err)
					continue
				}
			}
		}
	}
}

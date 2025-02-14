package collector

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/storage"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"math/rand"
	"runtime"
)

func Create(metrics *storage.MetricStorage) func() {
	pollCount := 0

	var stats runtime.MemStats

	return func() {
		runtime.ReadMemStats(&stats)

		pollCount++

		*metrics = storage.MetricStorage{
			"Alloc":         {Value: stats.Alloc, Type: consts.GAUGE},
			"BuckHashSys":   {Value: stats.BuckHashSys, Type: consts.GAUGE},
			"Frees":         {Value: stats.Frees, Type: consts.GAUGE},
			"GCCPUFraction": {Value: stats.GCCPUFraction, Type: consts.GAUGE},
			"GCSys":         {Value: stats.GCSys, Type: consts.GAUGE},
			"HeapAlloc":     {Value: stats.HeapAlloc, Type: consts.GAUGE},
			"HeapIdle":      {Value: stats.HeapIdle, Type: consts.GAUGE},
			"HeapInuse":     {Value: stats.HeapInuse, Type: consts.GAUGE},
			"HeapObjects":   {Value: stats.HeapObjects, Type: consts.GAUGE},
			"HeapReleased":  {Value: stats.HeapReleased, Type: consts.GAUGE},
			"HeapSys":       {Value: stats.HeapSys, Type: consts.GAUGE},
			"LastGC":        {Value: stats.LastGC, Type: consts.GAUGE},
			"Lookups":       {Value: stats.Lookups, Type: consts.GAUGE},
			"MCacheInuse":   {Value: stats.MCacheInuse, Type: consts.GAUGE},
			"MCacheSys":     {Value: stats.MCacheSys, Type: consts.GAUGE},
			"MSpanInuse":    {Value: stats.MSpanInuse, Type: consts.GAUGE},
			"MSpanSys":      {Value: stats.MSpanSys, Type: consts.GAUGE},
			"Mallocs":       {Value: stats.Mallocs, Type: consts.GAUGE},
			"NextGC":        {Value: stats.NextGC, Type: consts.GAUGE},
			"NumForcedGC":   {Value: stats.NumForcedGC, Type: consts.GAUGE},
			"NumGC":         {Value: stats.NumGC, Type: consts.GAUGE},
			"OtherSys":      {Value: stats.OtherSys, Type: consts.GAUGE},
			"PauseTotalNs":  {Value: stats.PauseTotalNs, Type: consts.GAUGE},
			"StackInuse":    {Value: stats.StackInuse, Type: consts.GAUGE},
			"StackSys":      {Value: stats.StackSys, Type: consts.GAUGE},
			"Sys":           {Value: stats.Sys, Type: consts.GAUGE},
			"TotalAlloc":    {Value: stats.TotalAlloc, Type: consts.GAUGE},
			"PollCount":     {Value: pollCount, Type: consts.COUNTER},
			"RandomValue":   {Value: rand.Float64(), Type: consts.GAUGE},
		}
		//logger.Instance.Info("CollectMetricsTask", "Данные", metrics)
	}
}

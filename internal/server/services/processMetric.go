package services

import (
	c "github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/files"
)

func (s *MetricService) ProcessMetric(metricData c.MetricData) (c.MetricData, error) {
	metricType := metricData.MType
	name := metricData.ID

	if metricType == c.GAUGE {
		s.store.SetGauge(name, metricData.Value)
	} else if metricType == c.COUNTER {
		s.store.SetCounter(name, metricData.Delta)
	}

	if s.storeInterval == 0 {
		filePath := s.fileStoragePath
		storeData := s.store.GetAll()
		files.WriteToFile(filePath, storeData)
	}

	newValue := s.store.Get(metricType, name)

	return c.NewMetricData(metricType, name, *newValue)
}

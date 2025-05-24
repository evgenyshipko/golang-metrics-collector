package services

import (
	"context"

	c "github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/files"
)

func (s *MetricService) ProcessMetric(ctx context.Context, metricData c.MetricData) (c.MetricData, error) {
	metricType := metricData.MType
	name := metricData.ID

	if metricType == c.GAUGE {
		s.store.SetGauge(ctx, name, metricData.Value)
	} else if metricType == c.COUNTER {
		s.store.SetCounter(ctx, name, metricData.Delta)
	}

	if s.storeInterval == 0 {
		filePath := s.fileStoragePath
		storeData, err := s.store.GetAll(ctx)
		if err != nil {
			return c.MetricData{}, err
		}

		files.WriteToFile(filePath, storeData)
	}

	newValue := s.store.Get(ctx, metricType, name)

	return c.NewMetricData(metricType, name, *newValue)
}

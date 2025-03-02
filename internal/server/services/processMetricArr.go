package services

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
)

func (s *MetricService) ProcessMetricArr(metricData []consts.MetricData) error {
	var storageData storage.StorageData
	for _, data := range metricData {
		storageData = append(storageData, storage.Data{
			Values: consts.Values{
				Gauge:   data.Value,
				Counter: data.Delta,
			},
			Name: data.ID,
		})
	}

	err := s.store.SetData(storageData)
	if err != nil {
		return err
	}

	return nil
}

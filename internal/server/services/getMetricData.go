package services

import (
	"errors"
	c "github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"net/http"
)

func (s *MetricService) GetMetricValue(metricData c.MetricData) (c.Values, error, int) {
	metricType := metricData.MType
	metricName := metricData.ID

	value := s.store.Get(metricType, metricName)
	if value.Counter == nil && value.Gauge == nil {
		return c.Values{}, errors.New("Метрики с таким именем нет в базе"), http.StatusNotFound
	}
	return *value, nil, 0
}

func (s *MetricService) GetMetricData(metricData c.MetricData) (c.MetricData, error, int) {
	metricType := metricData.MType
	metricName := metricData.ID

	value, err, status := s.GetMetricValue(metricData)
	if err != nil {
		return c.MetricData{}, err, status
	}
	data, err := c.NewMetricData(metricType, metricName, value)
	if err != nil {
		return c.MetricData{}, err, http.StatusBadRequest
	}
	return data, nil, 0
}

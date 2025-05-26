package services

import (
	"context"
	"errors"
	"net/http"

	c "github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
)

func (s *MetricService) GetMetricValue(ctx context.Context, metricData c.MetricData) (c.Values, int, error) {
	metricType := metricData.MType
	metricName := metricData.ID

	value := s.store.Get(ctx, metricType, metricName)
	if value.Counter == nil && value.Gauge == nil {
		return c.Values{}, http.StatusNotFound, errors.New("метрики с таким именем нет в базе")
	}
	return *value, 0, nil
}

func (s *MetricService) GetMetricData(ctx context.Context, metricData c.MetricData) (c.MetricData, int, error) {
	metricType := metricData.MType
	metricName := metricData.ID

	value, status, err := s.GetMetricValue(ctx, metricData)
	if err != nil {
		return c.MetricData{}, status, err
	}
	data, err := c.NewMetricData(metricType, metricName, value)
	if err != nil {
		return c.MetricData{}, http.StatusBadRequest, err
	}
	return data, 0, nil
}

package server

import (
	"errors"
	"fmt"
	c "github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/converter"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/url"
	"net/http"
)

func (s *Server) PostMetricOld(res http.ResponseWriter, req *http.Request) {
	metricType := c.Metric(url.MyURLParam(req, c.MetricType))
	name := url.MyURLParam(req, c.MetricName)
	value := req.Context().Value(c.MetricValue)

	logger.Instance.Infow("PostMetric", "metricType", metricType, "name", name, "value", value)

	err := s.store.Set(metricType, name, value)
	if err != nil {
		logger.Instance.Warn(fmt.Sprintf("PostMetric %s", errors.Unwrap(err)))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Write([]byte("Метрика записана успешно!"))
}

// TODO: покрыть тестами GET-хендлер
func (s *Server) GetMetricOld(res http.ResponseWriter, req *http.Request) {
	metricType := c.Metric(url.MyURLParam(req, c.MetricType))
	metricName := url.MyURLParam(req, c.MetricName)

	value := s.store.Get(metricType, metricName)
	if value == nil {
		http.Error(res, "Метрики с таким именем нет в базе", http.StatusNotFound)
		return
	}

	strVal, err := converter.MetricValueToString(metricType, value)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	logger.Instance.Infow("GetMetric", "metricType", metricType, "name", metricName, "value", strVal)

	res.Write([]byte(strVal))
}

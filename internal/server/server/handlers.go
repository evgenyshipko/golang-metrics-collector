package server

import (
	"encoding/json"
	"errors"
	"fmt"
	c "github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/converter"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/url"
	"net/http"
)

func (s *Server) PostMetric(res http.ResponseWriter, req *http.Request) {
	metricType := c.Metric(url.MyURLParam(req, c.MetricType))
	name := url.MyURLParam(req, c.MetricName)
	value := req.Context().Value(c.MetricValue)

	logger.Info("PostMetric", "metricType", metricType, "name", name, "value", value)

	err := s.store.Set(metricType, name, value)
	if err != nil {
		logger.Error(fmt.Sprintf("PostMetric %s", errors.Unwrap(err)))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Write([]byte("Метрика записана успешно!"))
}

// TODO: покрыть тестами GET-хендлер
func (s *Server) GetMetric(res http.ResponseWriter, req *http.Request) {
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

	logger.Info("GetMetric", "metricType", metricType, "name", metricName, "value", strVal)

	res.Write([]byte(strVal))
}

func (s *Server) NotFoundHandler(res http.ResponseWriter, _ *http.Request) {
	http.Error(res, "Запрашиваемый ресурс не найден", http.StatusNotFound)
}

func (s *Server) BadRequestHandler(res http.ResponseWriter, _ *http.Request) {
	http.Error(res, "URL не корректен", http.StatusBadRequest)
}

func (s *Server) ShowAllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	jsonStorage, err := json.MarshalIndent(s.store.GetAll(), "", "  ")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Info("ShowAllMetricsHandler", "jsonStorage", string(jsonStorage))
	data := fmt.Sprintf("<div>%s</div>", string(jsonStorage))
	res.Write([]byte(data))
}

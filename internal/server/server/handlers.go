package server

import (
	"encoding/json"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	m "github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares"
	"net/http"
)

func (s *CustomServer) StoreMetricHandler(res http.ResponseWriter, req *http.Request) {
	metricData, err := m.GetMetricDataFromContext(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	responseData, err := s.service.ProcessMetric(metricData)
	if err != nil {
		logger.Instance.Warnw("StoreMetricHandler", "GenerateMetricData", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	bytes, err := json.Marshal(responseData)
	if err != nil {
		logger.Instance.Warnw("StoreMetricHandler", "json.Marshal", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(bytes)
}

func (s *CustomServer) GetMetricDataHandler(res http.ResponseWriter, req *http.Request) {
	metricData, err := m.GetMetricDataFromContext(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	responseData, err, status := s.service.GetMetricData(metricData)
	if err != nil {
		logger.Instance.Warnw("GetMetricDataHandler", "GetMetricDataFromContext", err)
		http.Error(res, err.Error(), status)
		return
	}

	bytes, err := json.Marshal(responseData)
	if err != nil {
		logger.Instance.Warnw("GetMetricDataHandler", "json.Marshal", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(bytes)
}

func (s *CustomServer) GetMetricValueHandler(res http.ResponseWriter, req *http.Request) {
	metricData, err := m.GetMetricDataFromContext(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	value, err, status := s.service.GetMetricValue(metricData)
	if err != nil {
		logger.Instance.Warnw("GetMetricDataHandler", "GetMetricValue", err)
		http.Error(res, err.Error(), status)
		return
	}

	var strVal string
	if value.Counter != nil {
		strVal = fmt.Sprintf("%d", value.Counter)
	}
	if value.Gauge != nil {
		strVal = fmt.Sprintf("%d", value.Gauge)
	}

	res.Write([]byte(strVal))
}

func (s *CustomServer) ShowAllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	jsonStorage, err := json.MarshalIndent(*s.store.GetAll(), "", "  ")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Instance.Infow("ShowAllMetricsHandler", "jsonStorage", string(jsonStorage))
	data := fmt.Sprintf("<div>%s</div>", string(jsonStorage))
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.Write([]byte(data))
}

func (s *CustomServer) NotFoundHandler(res http.ResponseWriter, _ *http.Request) {
	http.Error(res, "Запрашиваемый ресурс не найден", http.StatusNotFound)
}

func (s *CustomServer) BadRequestHandler(res http.ResponseWriter, _ *http.Request) {
	http.Error(res, "URL не корректен", http.StatusBadRequest)
}

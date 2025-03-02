package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	m "github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares/update"
	middlewares "github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares/updates"
	"net/http"
	"strconv"
)

func (s *CustomServer) StoreMetricHandler(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), s.config.RequestWaitTimeout)
	defer cancel()

	metricData, err := m.GetMetricDataFromContext(ctx)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	responseData, err := s.service.ProcessMetric(ctx, metricData)
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

func (s *CustomServer) BatchStoreMetricHandler(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), s.config.RequestWaitTimeout)
	defer cancel()

	metricData, err := middlewares.GetArrayMetricDataFromContext(ctx)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.service.ProcessMetricArr(ctx, metricData)
	if err != nil {
		logger.Instance.Warnw("StoreMetricHandler", "ProcessMetricArr", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (s *CustomServer) GetMetricDataHandler(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), s.config.RequestWaitTimeout)
	defer cancel()

	metricData, err := m.GetMetricDataFromContext(ctx)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	responseData, status, err := s.service.GetMetricData(ctx, metricData)
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
	ctx, cancel := context.WithTimeout(req.Context(), s.config.RequestWaitTimeout)
	defer cancel()

	metricData, err := m.GetMetricDataFromContext(ctx)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	value, status, err := s.service.GetMetricValue(ctx, metricData)
	if err != nil {
		logger.Instance.Warnw("GetMetricDataHandler", "GetMetricValue", err)
		http.Error(res, err.Error(), status)
		return
	}

	var strVal string
	if value.Counter != nil {
		strVal = strconv.FormatInt(*value.Counter, 10)
	}
	if value.Gauge != nil {
		strVal = strconv.FormatFloat(*value.Gauge, 'f', -1, 64)
	}

	res.Write([]byte(strVal))
}

func (s *CustomServer) ShowAllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), s.config.RequestWaitTimeout)
	defer cancel()

	allMetrics, err := s.store.GetAll(ctx)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonStorage, err := json.MarshalIndent(*allMetrics, "", "  ")
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

func (s *CustomServer) PingDBConnection(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), s.config.RequestWaitTimeout)
	defer cancel()

	if !s.store.IsAvailable(ctx) {
		http.Error(res, "store not available", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

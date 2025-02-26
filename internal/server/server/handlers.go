package server

import (
	"encoding/json"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	m "github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares"
	"net/http"
	"strconv"
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

	responseData, status, err := s.service.GetMetricData(metricData)
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

	value, status, err := s.service.GetMetricValue(metricData)
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

func (s *CustomServer) PingDBConnection(res http.ResponseWriter, _ *http.Request) {
	dbPointer := s.GetDB()
	if dbPointer == nil {
		http.Error(res, "База данных не инициализирована", http.StatusInternalServerError)
		return
	}

	err := dbPointer.Ping()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

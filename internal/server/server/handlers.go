package server

import (
	"encoding/json"
	"fmt"
	c "github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/converter"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/files"
	m "github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares"
	"net/http"
)

func (s *CustomServer) StoreMetricHandler(res http.ResponseWriter, req *http.Request) {
	metricData, err := m.GetMetricData(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	metricType := metricData.MType
	name := metricData.ID

	if metricType == c.GAUGE {
		err = s.store.Set(metricType, name, *metricData.Value)
	} else if metricType == c.COUNTER {
		err = s.store.Set(metricType, name, *metricData.Delta)
	}

	if s.config.StoreInterval == 0 {
		filePath := s.config.FileStoragePath
		storeData := s.store.GetAll()
		files.WriteToFile(filePath, storeData)
	}

	if err != nil {
		logger.Instance.Warnw("StoreMetricHandler", "err in setting metric", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	newValue := s.store.Get(metricType, name)

	responseData, err := converter.GenerateMetricData(metricType, name, newValue)
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
	metricData, err := m.GetMetricData(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	metricType := metricData.MType
	metricName := metricData.ID

	value := s.store.Get(metricType, metricName)
	if value == nil {
		http.Error(res, "Метрики с таким именем нет в базе", http.StatusNotFound)
		return
	}

	responseData, err := converter.GenerateMetricData(metricType, metricName, value)
	if err != nil {
		logger.Instance.Warnw("GetMetricDataHandler", "GenerateMetricData", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	bytes, err := json.Marshal(responseData)
	if err != nil {
		logger.Instance.Warnw("GetMetricDataHandler", "json.Marshal", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Instance.Infow("GetMetricDataHandler", "metricType", metricType, "name", metricName, "value", value)

	res.Header().Set("Content-Type", "application/json")
	res.Write(bytes)
}

func (s *CustomServer) GetMetricValueHandler(res http.ResponseWriter, req *http.Request) {
	metricData, err := m.GetMetricData(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	metricType := metricData.MType
	metricName := metricData.ID

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

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	c "github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/converter"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	m "github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares"
	"net/http"
)

func (s *Server) StoreMetric(res http.ResponseWriter, req *http.Request) {
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

	if err != nil {
		logger.Instance.Warnw("err in setting metric", errors.Unwrap(err))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	newValue := s.store.Get(metricType, name)

	responseData, err := converter.GenerateMetricData(metricType, name, newValue)
	if err != nil {
		logger.Instance.Warnw(fmt.Sprintf("GenerateMetricData %s", errors.Unwrap(err)))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	bytes, err := json.Marshal(responseData)
	if err != nil {
		logger.Instance.Warnw(fmt.Sprintf("Marshal %s", errors.Unwrap(err)))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(bytes)
}

func (s *Server) GetMetric(res http.ResponseWriter, req *http.Request) {
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
		logger.Instance.Warnw(fmt.Sprintf("GenerateMetricData %s", errors.Unwrap(err)))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	bytes, err := json.Marshal(responseData)
	if err != nil {
		logger.Instance.Warnw(fmt.Sprintf("Marshal %s", errors.Unwrap(err)))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Instance.Infow("GetMetric", "metricType", metricType, "name", metricName, "value", value)

	res.Header().Set("Content-Type", "application/json")
	res.Write(bytes)
}

func (s *Server) ShowAllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	jsonStorage, err := json.MarshalIndent(s.store.GetAll(), "", "  ")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Instance.Infow("ShowAllMetricsHandler", "jsonStorage", string(jsonStorage))
	data := fmt.Sprintf("<div>%s</div>", string(jsonStorage))
	res.Write([]byte(data))
}

func (s *Server) NotFoundHandler(res http.ResponseWriter, _ *http.Request) {
	http.Error(res, "Запрашиваемый ресурс не найден", http.StatusNotFound)
}

func (s *Server) BadRequestHandler(res http.ResponseWriter, _ *http.Request) {
	http.Error(res, "URL не корректен", http.StatusBadRequest)
}

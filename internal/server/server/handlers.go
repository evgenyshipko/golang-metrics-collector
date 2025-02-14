package server

import (
	"encoding/json"
	"errors"
	"fmt"
	c "github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
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

	responseData := generateMetricData(metricType, name, newValue)

	bytes, err := json.Marshal(responseData)
	if err != nil {
		logger.Instance.Warnw(fmt.Sprintf("StoreMetric %s", errors.Unwrap(err)))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

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

	responseData := generateMetricData(metricType, metricName, value)

	bytes, err := json.Marshal(responseData)
	if err != nil {
		logger.Instance.Warnw(fmt.Sprintf("StoreMetric %s", errors.Unwrap(err)))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Instance.Infow("GetMetric", "metricType", metricType, "name", metricName, "value", value)

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

func generateMetricData(metricType c.Metric, name string, value interface{}) m.MetricData {
	var GaugeValue *float64
	var CounterValue *int64
	if metricType == c.COUNTER {
		if v, ok := value.(int64); ok {
			CounterValue = &v
		}
	} else if metricType == c.GAUGE {
		if v, ok := value.(float64); ok {
			GaugeValue = &v
		}
	}

	return m.MetricData{
		ID:    name,
		MType: metricType,
		Value: GaugeValue,
		Delta: CounterValue,
	}
}

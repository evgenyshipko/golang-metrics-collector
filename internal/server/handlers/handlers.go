package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/converter"
	"github.com/evgenyshipko/golang-metrics-collector/internal/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
	"github.com/go-chi/chi"
	"net/http"
)

func PostMetric(res http.ResponseWriter, req *http.Request) {
	metricType := consts.Metric(chi.URLParam(req, "metricType"))
	name := chi.URLParam(req, "metricName")
	value := req.Context().Value("metricValue")

	err := storage.STORAGE.Set(metricType, name, value)
	if err != nil {
		logger.Error(fmt.Sprintf("PostMetric %s", errors.Unwrap(err)))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info("PostMetric", "Data after save", storage.STORAGE.GetAll())

	res.Write([]byte("Метрика записана успешно!"))
}

// TODO: покрыть тестами GET-хендлер
func GetMetric(res http.ResponseWriter, req *http.Request) {
	metricType := consts.Metric(chi.URLParam(req, "metricType"))
	metricName := chi.URLParam(req, "metricName")

	value := storage.STORAGE.Get(metricType, metricName)
	if value == nil {
		http.Error(res, "Метрики с таким именем нет в базе", http.StatusNotFound)
		return
	}

	strVal, err := converter.MetricValueToString(metricName, value)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	res.Write([]byte(strVal))
}

func NotFoundHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Запрашиваемый ресурс не найден", http.StatusNotFound)
	return
}

func BadRequestHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "URL не корректен", http.StatusBadRequest)
	return
}

func ShowAllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	jsonStorage, err := json.MarshalIndent(storage.STORAGE.GetAll(), "", "  ")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Info("ShowAllMetricsHandler", "jsonStorage", string(jsonStorage))
	data := fmt.Sprintf("<div>%s</div>", string(jsonStorage))
	res.Write([]byte(data))
}

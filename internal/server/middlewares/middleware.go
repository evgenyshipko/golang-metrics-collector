package middlewares

import (
	"context"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/consts"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

func ValidateMetricType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricType := consts.Metric(chi.URLParam(r, "metricType"))
		if metricType != consts.COUNTER && metricType != consts.GAUGE {
			http.Error(w, "Неизвестный тип метрики", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ValidateMetricValue(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricValue := chi.URLParam(r, "metricValue")

		metricType := consts.Metric(chi.URLParam(r, "metricType"))

		ctx := context.WithValue(r.Context(), "metricType", metricType)

		if metricValue != "" {

			if metricType == consts.GAUGE {
				float64Value, err := strconv.ParseFloat(metricValue, 64)
				if err != nil {
					http.Error(w, "неверное Value для gauge", http.StatusBadRequest)
					return
				}

				ctx = context.WithValue(r.Context(), "metricValue", float64Value)

			} else if metricType == consts.COUNTER {
				int64Value, err := strconv.ParseInt(metricValue, 10, 64)
				if err != nil {
					http.Error(w, "неверное Value для counter", http.StatusBadRequest)
					return
				}

				ctx = context.WithValue(r.Context(), "metricValue", int64Value)

			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

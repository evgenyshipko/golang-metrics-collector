package middlewares

import (
	"context"
	c "github.com/evgenyshipko/golang-metrics-collector/internal/server/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/url"
	"net/http"
	"strconv"
)

func ValidateMetricType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricType := c.Metric(url.URLParam(r, c.MetricType))
		if metricType != c.COUNTER && metricType != c.GAUGE {
			http.Error(w, "Неизвестный тип метрики", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ValidateMetricValue(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricValue := url.URLParam(r, c.MetricValue)

		metricType := c.Metric(url.URLParam(r, c.MetricType))

		ctx := context.WithValue(r.Context(), c.MetricType, metricType)

		if metricValue != "" {

			if metricType == c.GAUGE {
				float64Value, err := strconv.ParseFloat(metricValue, 64)
				if err != nil {
					http.Error(w, "неверное Value для gauge", http.StatusBadRequest)
					return
				}

				ctx = context.WithValue(r.Context(), c.MetricValue, float64Value)

			} else if metricType == c.COUNTER {
				int64Value, err := strconv.ParseInt(metricValue, 10, 64)
				if err != nil {
					http.Error(w, "неверное Value для counter", http.StatusBadRequest)
					return
				}

				ctx = context.WithValue(r.Context(), c.MetricValue, int64Value)

			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

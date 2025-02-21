package middlewares

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"net/http"
)

func ValidateType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricData, err := GetMetricData(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		metricType := metricData.MType

		if metricType != consts.COUNTER && metricType != consts.GAUGE {
			http.Error(w, "Неизвестный тип метрики", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

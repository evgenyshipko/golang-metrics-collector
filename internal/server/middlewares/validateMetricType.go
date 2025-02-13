package middlewares

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/url"
	"net/http"
)

func ValidateMetricType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricType := consts.Metric(url.MyURLParam(r, consts.MetricType))
		if metricType != consts.COUNTER && metricType != consts.GAUGE {
			http.Error(w, "Неизвестный тип метрики", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

package middlewares

import (
	c "github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"net/http"
)

func ValidateValue(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		metricData, err := GetMetricDataFromContext(req.Context())
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		logger.Instance.Debugw("ValidateValue", "metricData", metricData)

		if (metricData.Delta != nil && metricData.Value != nil) || (metricData.Delta == nil && metricData.Value == nil) {
			http.Error(res, "Value метрики может храниться либо в Delta либо в Value", http.StatusBadRequest)
			return
		}

		metricType := metricData.MType

		if metricType == c.GAUGE && metricData.Value == nil {
			http.Error(res, "отсутствует Value для gauge", http.StatusBadRequest)
			return
		} else if metricType == c.COUNTER && metricData.Delta == nil {
			http.Error(res, "отсутствует Value для counter", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(res, req)
	})
}

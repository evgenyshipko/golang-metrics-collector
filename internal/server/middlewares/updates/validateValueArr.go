package middlewares

import (
	"net/http"

	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/validate"
)

func ValidateValueArr(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		metricData, err := GetArrayMetricDataFromContext(req.Context())
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		logger.Instance.Debugw("ValidateValue", "metricData", metricData)

		for _, metric := range metricData {
			err = validate.Value(metric)
			if err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
		}

		next.ServeHTTP(res, req)
	})
}

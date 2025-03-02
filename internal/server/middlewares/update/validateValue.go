package update

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/validate"
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

		err = validate.Value(metricData)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		next.ServeHTTP(res, req)
	})
}

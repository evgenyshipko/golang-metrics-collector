package middlewares

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/validate"
	"net/http"
)

func ValidateTypeArr(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricData, err := GetArrayMetricDataFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, metric := range metricData {
			err = validate.Type(metric)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

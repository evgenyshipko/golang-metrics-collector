package middlewares

import (
	"net/http"

	"github.com/evgenyshipko/golang-metrics-collector/internal/server/validate"
)

func ValidateNameArr(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricData, err := GetArrayMetricDataFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, metric := range metricData {
			err = validate.Name(metric)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

package update

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/validate"
	"net/http"
)

func ValidateType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricData, err := GetMetricDataFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = validate.Type(metricData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

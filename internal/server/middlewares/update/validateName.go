package update

import (
	"net/http"

	"github.com/evgenyshipko/golang-metrics-collector/internal/server/validate"
)

func ValidateName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricData, err := GetMetricDataFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = validate.Name(metricData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

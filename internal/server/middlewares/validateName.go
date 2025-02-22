package middlewares

import (
	"net/http"
)

func ValidateName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricData, err := GetMetricDataFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if metricData.ID == "" {
			http.Error(w, "Не было передано имя метрики", http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

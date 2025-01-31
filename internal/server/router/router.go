package router

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/handlers"
	m "github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares"
	"github.com/go-chi/chi"
)

func MakeChiRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", handlers.ShowAllMetricsHandler)
	r.Get("/value/{metricType}/{metricName}", handlers.GetMetric)
	r.Post("/update/", handlers.BadRequestHandler)

	r.With(m.ValidateMetricType).Route("/update/{metricType}", func(r chi.Router) {
		r.Post("/", handlers.NotFoundHandler)

		r.With(m.ValidateMetricValue).Post("/{metricValue}", handlers.NotFoundHandler)

		r.With(m.ValidateMetricValue).Post("/{metricName}/{metricValue}", handlers.PostMetric)
	})
	return r
}

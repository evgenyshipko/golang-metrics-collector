package server

import (
	m "github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares"
	"github.com/go-chi/chi"
)

func (s *Server) routes() {
	s.router.Get("/", s.ShowAllMetricsHandler)
	s.router.Get("/value/{metricType}/{metricName}", s.GetMetric)
	s.router.Post("/update/", s.BadRequestHandler)

	s.router.With(m.ValidateMetricType).Route("/update/{metricType}", func(r chi.Router) {
		r.Post("/", s.NotFoundHandler)

		r.With(m.ValidateMetricValue).Post("/{metricValue}", s.NotFoundHandler)

		r.With(m.ValidateMetricValue).Post("/{metricName}/{metricValue}", s.PostMetric)
	})
}

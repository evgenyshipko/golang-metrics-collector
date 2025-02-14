package server

import (
	m "github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares"
	"github.com/go-chi/chi"
)

func (s *Server) routes() {
	s.router.Get("/", s.ShowAllMetricsHandler)

	// TODO: написать тесты на ручку
	s.router.Get("/value/{metricType}/{metricName}", s.GetMetricOld)

	s.router.With(m.SaveBodyToContext, m.ValidateName, m.ValidateType).Route("/value", func(r chi.Router) {
		r.Post("/", s.GetMetric)
	})

	s.router.With(m.SaveBodyToContext, m.ValidateName, m.ValidateType, m.ValidateValue).Route("/update", func(r chi.Router) {
		r.Post("/", s.StoreMetric)
	})

	s.router.With(m.ValidateTypeByURLParam).Route("/update/{metricType}", func(r chi.Router) {
		r.Post("/", s.NotFoundHandler)

		r.With(m.ValidateValueByURLParam).Post("/{metricValue}", s.NotFoundHandler)

		r.With(m.ValidateValueByURLParam).Post("/{metricName}/{metricValue}", s.PostMetricOld)
	})
}

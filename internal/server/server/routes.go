package server

import (
	m "github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares"
	"github.com/go-chi/chi"
)

func (s *Server) routes() {
	s.router.Get("/", s.ShowAllMetricsHandler)

	s.router.Route("/value", func(r chi.Router) {
		r.With(m.SaveBodyToContext, m.ValidateName, m.ValidateType).Post("/", s.GetMetric)

		r.With(m.SaveURLParamsToContext, m.ValidateName, m.ValidateType).Get("/{metricType}/{metricName}", s.GetMetric)
	})

	s.router.Route("/update", func(r chi.Router) {
		r.With(m.SaveBodyToContext, m.ValidateName, m.ValidateType, m.ValidateValue).Post("/", s.StoreMetric)

		r.With(m.SaveURLParamsToContext, m.ValidateType).Post("/{metricType}", s.NotFoundHandler)

		r.With(m.SaveURLParamsToContext, m.ValidateType, m.ValidateValue).Post("/{metricType}/{metricValue}", s.NotFoundHandler)

		r.With(m.SaveURLParamsToContext, m.ValidateName, m.ValidateType, m.ValidateValue).Post("/{metricType}/{metricName}/{metricValue}", s.StoreMetric)
	})
}

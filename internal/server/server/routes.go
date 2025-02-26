package server

import (
	m "github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares"
	"github.com/go-chi/chi"
)

func (s *CustomServer) routes() {
	s.router.Get("/", s.ShowAllMetricsHandler)

	s.router.Get("/ping", s.PingDBConnection)

	s.router.Route("/value", func(r chi.Router) {
		r.With(m.SaveBodyToContext, m.ValidateName, m.ValidateType).Post("/", s.GetMetricDataHandler)

		r.With(m.SaveURLParamsToContext, m.ValidateName, m.ValidateType).Get("/{metricType}/{metricName}", s.GetMetricValueHandler)
	})

	s.router.Route("/update", func(r chi.Router) {
		r.With(m.SaveBodyToContext, m.ValidateName, m.ValidateType, m.ValidateValue).Post("/", s.StoreMetricHandler)

		r.With(m.SaveURLParamsToContext, m.ValidateName, m.ValidateType, m.ValidateValue).Post("/{metricType}/{metricName}/{metricValue}", s.StoreMetricHandler)
	})
}

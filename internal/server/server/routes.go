package server

import (
	m "github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares"
	"github.com/go-chi/chi"
)

func (s *Server) routes() {
	s.router.Get("/", s.ShowAllMetricsHandler)

	// TODO: написать тесты на ручку
	s.router.Get("/value/{metricType}/{metricName}", s.GetMetricOld)

	//s.router.With(m.SaveBodyToContext, m.ValidateName, m.ValidateType).Route("/value", func(r chi.Router) {
	//	r.Post("/", s.GetMetric)
	//})

	//s.router.With(m.SaveBodyToContext, m.ValidateName, m.ValidateType, m.ValidateValue).Route("/update", func(r chi.Router) {
	//	r.Post("/", s.StoreMetric)
	//})

	s.router.Route("/update", func(r chi.Router) {
		r.With(m.SaveBodyToContext, m.ValidateName, m.ValidateType, m.ValidateValue).Post("/", s.StoreMetric)

		r.With(m.SaveURLParamsToContext, m.ValidateType).Post("/{metricType}", s.NotFoundHandler)

		r.With(m.SaveURLParamsToContext, m.ValidateType, m.ValidateValue).Post("/{metricType}/{metricValue}", s.NotFoundHandler)

		r.With(m.SaveURLParamsToContext, m.ValidateName, m.ValidateType, m.ValidateValue).Post("/{metricType}/{metricName}/{metricValue}", s.StoreMetric)
	})

	//s.router.With(m.SaveURLParamsToContext, m.ValidateName, m.ValidateType, m.ValidateValue).Route("/update/{metricType}/{metricName}/{metricValue}", func(r chi.Router) {
	//	r.Post("/", s.StoreMetric)
	//})
}

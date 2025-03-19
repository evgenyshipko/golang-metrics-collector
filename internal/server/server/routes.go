package server

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares/update"
	middlewares "github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares/updates"
	"github.com/go-chi/chi"
)

func (s *CustomServer) routes() {
	s.router.Get("/", s.ShowAllMetricsHandler)

	s.router.Get("/ping", s.PingDBConnection)

	s.router.Route("/value", func(r chi.Router) {
		r.With(update.SaveBodyToContext, update.ValidateName, update.ValidateType).Post("/", s.GetMetricDataHandler)

		r.With(update.SaveURLParamsToContext, update.ValidateName, update.ValidateType).Get("/{metricType}/{metricName}", s.GetMetricValueHandler)
	})

	s.router.Route("/update", func(r chi.Router) {
		r.With(update.SaveBodyToContext, update.ValidateName, update.ValidateType, update.ValidateValue).Post("/", s.StoreMetricHandler)

		r.With(update.SaveURLParamsToContext, update.ValidateName, update.ValidateType, update.ValidateValue).Post("/{metricType}/{metricName}/{metricValue}", s.StoreMetricHandler)
	})

	s.router.With(middlewares.SaveMetricsToContext, middlewares.ValidateNameArr, middlewares.ValidateTypeArr, middlewares.ValidateValueArr).Post("/updates/", s.BatchStoreMetricHandler)
}

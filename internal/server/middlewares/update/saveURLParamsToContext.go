package update

import (
	"context"
	c "github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/url"
	"net/http"
	"strconv"
)

func SaveURLParamsToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Instance.Info("SaveURLParamsToContext")

		metricValue := url.MyURLParam(r, c.MetricValue)

		metricType := c.Metric(url.MyURLParam(r, c.MetricType))

		metricName := url.MyURLParam(r, c.MetricName)

		metricData := c.MetricData{
			ID:    metricName,
			MType: metricType,
		}

		if metricValue != "" {

			if metricType == c.GAUGE {
				float64Value, err := strconv.ParseFloat(metricValue, 64)
				if err != nil {
					http.Error(w, "неверное Value для gauge", http.StatusBadRequest)
					return
				}

				metricData.Value = &float64Value

			} else if metricType == c.COUNTER {
				int64Value, err := strconv.ParseInt(metricValue, 10, 64)
				if err != nil {
					http.Error(w, "неверное Value для counter", http.StatusBadRequest)
					return
				}

				metricData.Delta = &int64Value

			}
		}

		logger.Instance.Debugw("SaveURLParamsToContext", "metricData", metricData)

		ctx := context.WithValue(r.Context(), MetricDataKey, metricData)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

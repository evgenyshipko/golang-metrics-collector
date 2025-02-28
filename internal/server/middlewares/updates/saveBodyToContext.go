package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	c "github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"net/http"
)

type ContextKey string

const MetricDataArrayKey ContextKey = "metricDataArray"

func SaveBodyArrayToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		var metricData []c.MetricData
		var buf bytes.Buffer
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			logger.Instance.Warnw(err.Error())
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &metricData); err != nil {
			logger.Instance.Warnw(err.Error())
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		logger.Instance.Debugw("SaveBodyArrayToContext", "metricData", metricData)

		ctx := context.WithValue(req.Context(), MetricDataArrayKey, metricData)

		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

func GetArrayMetricDataFromContext(ctx context.Context) ([]c.MetricData, error) {
	metricData := ctx.Value(MetricDataArrayKey)

	data, ok := metricData.([]c.MetricData)
	if !ok {
		logger.Instance.Warn("Невозможно привести к []MetricData")
		return []c.MetricData{}, errors.New("невозможно привести к []MetricData")
	}
	return data, nil
}

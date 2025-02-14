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

const metricDataKey ContextKey = "metricData"

func SaveBodyToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		var metricData c.MetricData
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

		ctx := context.WithValue(req.Context(), metricDataKey, metricData)

		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

func GetMetricData(ctx context.Context) (c.MetricData, error) {
	metricData := ctx.Value(metricDataKey)

	data, ok := metricData.(c.MetricData)
	if !ok {
		logger.Instance.Warnw("Невозможно привести к MetricData")
		return c.MetricData{}, errors.New("невозможно привести к MetricData")
	}
	return data, nil
}

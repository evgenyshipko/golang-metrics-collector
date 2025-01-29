package helpers

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/convert"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/consts"
	"strconv"
)

func GetMetricType(metricName string) consts.Metric {
	if metricName == "PollCount" {
		return consts.COUNTER
	}
	return consts.GAUGE
}

func ConvertMetricValueToString(name string, value interface{}) (string, error) {
	stringValue := ""
	metricType := GetMetricType(name)
	if metricType == consts.GAUGE {
		float64val, err := convert.ToFloat64(value)
		if err != nil {
			return "", err
		}
		stringValue = strconv.FormatFloat(float64val, 'f', 6, 64)
	}

	if metricType == consts.COUNTER {
		int64val, err := convert.ToInt64(value)
		if err != nil {
			return "", err
		}
		stringValue = strconv.FormatInt(int64val, 10)
	}
	return stringValue, nil
}

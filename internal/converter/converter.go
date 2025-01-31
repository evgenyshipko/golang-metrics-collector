package converter

import (
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/helpers"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/consts"
	"strconv"
)

func ToInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	default:
		return 0, fmt.Errorf("%T не конвертируется в int64", v)
	}
}

// при конвертации int64 и uint64 теряется точность
func ToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("%T не конвертируется в float64", v)
	}
}

// TODO: вынести metricType отдельным параметром
func MetricValueToString(name string, value interface{}) (string, error) {
	stringValue := ""
	metricType := helpers.GetMetricType(name)
	if metricType == consts.GAUGE {
		float64val, err := ToFloat64(value)
		if err != nil {
			return "", fmt.Errorf("ошибка в MetricValueToString, metricType: %s, %w", metricType, err)
		}
		stringValue = strconv.FormatFloat(float64val, 'f', 6, 64)
	}

	if metricType == consts.COUNTER {
		int64val, err := ToInt64(value)
		if err != nil {
			return "", fmt.Errorf("ошибка в MetricValueToString, metricType: %s, %w", metricType, err)
		}
		stringValue = strconv.FormatInt(int64val, 10)
	}
	return stringValue, nil
}

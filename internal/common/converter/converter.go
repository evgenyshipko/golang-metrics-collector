package converter

import (
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"reflect"
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

// помнить, что при конвертации int64 и uint64 теряется точность
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

func MetricValueToString(metricType consts.Metric, value interface{}) (string, error) {
	stringValue := ""
	if metricType == consts.GAUGE {
		float64val, err := ToFloat64(value)
		if err != nil {
			return "", fmt.Errorf("ошибка в MetricValueToString, metricType: %s, %w", metricType, err)
		}
		stringValue = strconv.FormatFloat(float64val, 'f', -1, 64)
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

func GenerateMetricData(metricType consts.Metric, name string, value interface{}) (consts.MetricData, error) {
	var GaugeValue *float64
	var CounterValue *int64
	logger.Instance.Debugw("GenerateMetricData", "value", value, "type", reflect.TypeOf(value).String())
	if metricType == consts.COUNTER {
		intVal, err := ToInt64(value)
		if err != nil {
			return consts.MetricData{}, err
		}
		CounterValue = &intVal
	} else if metricType == consts.GAUGE {
		floatVal, err := ToFloat64(value)
		if err != nil {
			return consts.MetricData{}, err
		}
		GaugeValue = &floatVal
	}

	return consts.MetricData{
		ID:    name,
		MType: metricType,
		Value: GaugeValue,
		Delta: CounterValue,
	}, nil
}

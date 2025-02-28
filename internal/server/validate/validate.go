package validate

import (
	"errors"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
)

func Name(metricData consts.MetricData) error {
	if metricData.ID == "" {
		return errors.New("не было передано имя метрики")
	}
	return nil
}

func Type(metricData consts.MetricData) error {
	metricType := metricData.MType

	if metricType != consts.COUNTER && metricType != consts.GAUGE {
		return errors.New("неизвестный тип метрики")
	}
	return nil
}

func Value(metricData consts.MetricData) error {
	if (metricData.Delta != nil && metricData.Value != nil) || (metricData.Delta == nil && metricData.Value == nil) {
		return errors.New("value метрики может храниться либо в Delta либо в Value")
	}

	metricType := metricData.MType

	if metricType == consts.GAUGE && metricData.Value == nil {
		return errors.New("отсутствует Value для gauge")
	} else if metricType == consts.COUNTER && metricData.Delta == nil {
		return errors.New("отсутствует Value для counter")
	}
	return nil
}

package parser

import (
	"errors"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/consts"
	"regexp"
	"strconv"
)

type MetricData struct {
	MetricType consts.Metric
	Name       string
	Value      interface{}
}

func IsNameMissed(path string) bool {
	re := regexp.MustCompile(`^/update/(gauge|counter)/(-?\d+(\.\d+)?)?$`)

	matches := re.FindStringSubmatch(path)
	return len(matches) > 0
}

func ParseURLPath(path string) (MetricData, error) {

	re := regexp.MustCompile(`^/update/(gauge|counter)/([^/]+)/(-?\d+(\.\d+)?)$`)

	matches := re.FindStringSubmatch(path)
	if len(matches) == 0 {
		return MetricData{}, errors.New("неверный формат URL")
	}

	metricType := consts.Metric(matches[1])
	metricName := matches[2]
	metricValueStr := matches[3]

	if metricType == consts.GAUGE {
		floatVal, err := strconv.ParseFloat(metricValueStr, 64)
		if err != nil {
			return MetricData{}, errors.New("неверное Value для gauge")
		}
		return MetricData{metricType, metricName, floatVal}, nil
	} else if metricType == consts.COUNTER {
		intVal, err := strconv.ParseInt(metricValueStr, 10, 64)
		if err != nil {
			return MetricData{}, errors.New("неверное Value для counter")
		}
		return MetricData{metricType, metricName, intVal}, nil
	}

	return MetricData{}, errors.New("не может быть, но все же как-то программа до сюда дошла :D")
}

package storage

import "github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"

type Data struct {
	consts.Values
	Name string `json:"name"`
}

type StorageData []Data

type Storage interface {
	Get(metricType consts.Metric, name string) *consts.Values
	GetAll() (*StorageData, error)
	SetGauge(name string, value *float64)
	SetCounter(name string, value *int64)
	SetData(data StorageData) error
}

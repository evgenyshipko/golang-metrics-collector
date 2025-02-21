package storage

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
)

type Data struct {
	Name    string   `json:"name"`
	Counter *int64   `json:"counter,omitempty"`
	Gauge   *float64 `json:"gauge,omitempty"`
}

type MemStorageData []Data

type MemStorage struct {
	data MemStorageData
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: MemStorageData{},
	}
}

func (storage *MemStorage) Get(metricType consts.Metric, name string) interface{} {
	if metricType == consts.COUNTER {
		dataPointer := storage.getCounterData(name)
		if dataPointer != nil && dataPointer.Counter != nil {
			return *dataPointer.Counter
		}
	}
	if metricType == consts.GAUGE {
		dataPointer := storage.getGaugeData(name)
		if dataPointer != nil && dataPointer.Gauge != nil {
			return *dataPointer.Gauge
		}
	}
	return nil
}

func (storage *MemStorage) getCounterData(name string) *Data {
	for index, data := range storage.data {
		if data.Name == name && data.Counter != nil {
			return &storage.data[index]
		}
	}
	return nil
}

func (storage *MemStorage) getGaugeData(name string) *Data {
	for index, data := range storage.data {
		if data.Name == name && data.Gauge != nil {
			// ЗАПОМНИТЬ: берем значение из слайса по индексу т.к. нам далее нужно менять оригинальне значение,
			// а data - это копия и менять ее поля не имеет смысла
			return &storage.data[index]
		}
	}
	return nil
}

func (storage *MemStorage) SetGauge(name string, value *float64) {
	dataPointer := storage.getGaugeData(name)
	if dataPointer != nil {
		dataPointer.Gauge = value
		return
	}
	storage.data = append(storage.data, Data{Name: name, Gauge: value})
}

func (storage *MemStorage) SetCounter(name string, value *int64) {
	dataPointer := storage.getCounterData(name)
	if dataPointer != nil {
		existingValue := *dataPointer.Counter
		resultValue := existingValue + *value
		dataPointer.Counter = &resultValue
		return
	}
	storage.data = append(storage.data, Data{Name: name, Counter: value})
}

func (storage *MemStorage) GetAll() *MemStorageData {
	return &storage.data
}

func (storage *MemStorage) SetData(data MemStorageData) {
	storage.data = data
}

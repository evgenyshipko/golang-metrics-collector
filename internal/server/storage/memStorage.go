package storage

import (
	"context"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
)

type MemStorage struct {
	data StorageData
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: StorageData{},
	}
}

func (storage *MemStorage) Get(_ context.Context, metricType consts.Metric, name string) *consts.Values {
	if metricType == consts.COUNTER {
		dataPointer := storage.getCounterData(name)
		if dataPointer != nil && dataPointer.Counter != nil {
			return &consts.Values{
				Counter: dataPointer.Counter,
			}
		}
	}
	if metricType == consts.GAUGE {
		dataPointer := storage.getGaugeData(name)
		if dataPointer != nil && dataPointer.Gauge != nil {
			return &consts.Values{
				Gauge: dataPointer.Gauge,
			}
		}
	}
	return &consts.Values{}
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
			// ЗАПОМНИТЬ: берем значение из слайса по индексу т.к. нам далее нужно менять оригинальное значение,
			// а data - это копия и менять ее поля не имеет смысла
			return &storage.data[index]
		}
	}
	return nil
}

func (storage *MemStorage) SetGauge(_ context.Context, name string, value *float64) {
	dataPointer := storage.getGaugeData(name)
	if dataPointer != nil {
		dataPointer.Gauge = value
		return
	}
	storage.data = append(storage.data, Data{Name: name, Values: consts.Values{Gauge: value}})
}

func (storage *MemStorage) SetCounter(_ context.Context, name string, value *int64) {
	dataPointer := storage.getCounterData(name)
	if dataPointer != nil {
		existingValue := *dataPointer.Counter
		resultValue := existingValue + *value
		dataPointer.Counter = &resultValue
		return
	}
	storage.data = append(storage.data, Data{Name: name, Values: consts.Values{Counter: value}})
}

func (storage *MemStorage) GetAll(_ context.Context) (*StorageData, error) {
	return &storage.data, nil
}

func (storage *MemStorage) SetData(ctx context.Context, storageData StorageData) error {
	for _, data := range storageData {
		if data.Counter != nil {
			storage.SetCounter(ctx, data.Name, data.Counter)
		}
		if data.Gauge != nil {
			storage.SetGauge(ctx, data.Name, data.Gauge)
		}
	}
	return nil
}

func (storage *MemStorage) IsAvailable(_ context.Context) bool {
	return true
}

func (storage *MemStorage) Close() error {
	return nil
}

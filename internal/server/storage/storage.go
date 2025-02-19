package storage

import (
	"errors"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/converter"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"reflect"
	"strings"
)

type MemStorageData map[string]interface{}

type MemStorage struct {
	data MemStorageData
}

type MemStorageInterface[V comparable] interface {
	Get(metricType string, name string) V
	Set(metricType string, name string, value V) error
	GetAll() MemStorageData
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: make(MemStorageData),
	}
}

func getKey(metricType consts.Metric, metricName string) string {
	return fmt.Sprintf(`%s_%s`, metricType, metricName)
}

func (storage *MemStorage) Get(metricType consts.Metric, name string) interface{} {
	key := getKey(metricType, name)
	return storage.data[key]
}

func (storage *MemStorage) Set(metricType consts.Metric, name string, value interface{}) error {

	logger.Instance.Debugw("MemStorage.Set", "type", metricType, "name", name, "value", value, "type", reflect.TypeOf(value).String())

	key := getKey(metricType, name)

	if metricType == consts.COUNTER {

		int64Value, err := converter.ToInt64(value)
		if err != nil {
			return fmt.Errorf("ошибка в Set, metricType: %s, %w", metricType, err)
		}

		if storage.data[key] != nil {
			prevInt64Value, err := converter.ToInt64(storage.data[key])
			if err != nil {
				return fmt.Errorf("ошибка в Set, metricType: %s, %w", metricType, err)
			}
			storage.data[key] = prevInt64Value + int64Value
		} else {
			storage.data[key] = int64Value
		}
		return nil
	}

	if metricType == consts.GAUGE {

		float64Value, err := converter.ToFloat64(value)
		if err != nil {
			return fmt.Errorf("ошибка в Set, metricType: %s, %w", metricType, err)
		}

		storage.data[key] = float64Value
		return nil
	}

	return errors.New("неизвестный тип метрики")
}

func (storage *MemStorage) GetAll() *MemStorageData {
	return &storage.data
}

func (storage *MemStorage) SetData(data MemStorageData) {

	// восстанавливаем правильный тип для counter
	for key, value := range data {
		if strings.HasPrefix(key, "counter_") {
			if floatVal, ok := value.(float64); ok {
				data[key] = int64(floatVal)
			}
		}
	}

	storage.data = data
}

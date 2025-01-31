package storage

import (
	"errors"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/converter"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/consts"
)

type MemStorage struct {
	data map[string]interface{}
}

// TODO: это скорее MetricMemStorageInterface т.к. есть зависимость от имении метрики. Сделать упрощенный MemStorageInterface
type MemStorageInterface[V comparable] interface {
	Get(metricType string, name string) V
	Set(metricType string, name string, value V) error
	GetAll() map[string]interface{}
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: make(map[string]interface{}),
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

	//logger.Info(string(metricType), "name", name, "value", value, "type", reflect.TypeOf(value).String())

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
			/* ВОПРОС РЕВЬЮВЕРУ:
			Если типизировать следующим образом: storage.data[key].(consts.Counter) + value.(consts.Counter),
			То получаем панику:  interface conversion: interface {} is int64, not consts.Counter
			Хотя type Counter - это int64. Странная штуковина, я хочу обозвать именованным типом, а программа ругается.
			Это ограничение языка или я что-то не так делаю?
			*/
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

func (storage *MemStorage) GetAll() map[string]interface{} {
	return storage.data
}

var STORAGE = NewMemStorage()

package storage

import (
	"errors"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/convert"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/consts"
	"reflect"
)

type MemStorage struct {
	data map[string]interface{}
}

type MemStorageInterface[V comparable] interface {
	Get(metricType string, name string) (V, error)
	Set(metricType string, name string, value V) error
	Delete(metricType string, name string) error
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: make(map[string]interface{}),
	}
}

/*
ВОПРОС РЕВЬЮВЕРУ:
Насколько правильно зашивать логику сохранения значения в зависимости от типа метрики в эту функцию?
*/

func (storage *MemStorage) Set(metricType consts.Metrics, name string, value interface{}) error {

	fmt.Println(metricType, name, value, reflect.TypeOf(value))

	key := fmt.Sprintf(`%s_%s`, metricType, name)

	if metricType == consts.COUNTER {

		int64Value, err := convert.ToInt64(value)
		if err != nil {
			return err
		}

		if storage.data[key] != nil {
			prevInt64Value, err := convert.ToInt64(value)
			if err != nil {
				return err
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

		float64Value, err := convert.ToFloat64(value)
		if err != nil {
			return err
		}

		storage.data[key] = float64Value
		return nil
	}

	return errors.New("неизвестный тип метрики")
}

var STORAGE = NewMemStorage()

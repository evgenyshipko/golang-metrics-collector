package storage

import (
	"errors"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/consts"
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
ВОПРОС ПРОВЕРЯЮЩЕМУ
Как нормально типизировать value? Мы точно знаем из ТЗ, что у gauge тип float64, а у counter - int64.
Но сколько я не бился - не смог понять как эти знания можно применить в типизации функции/стора.
Пробовал дженерики - не подходят, т.к при их использовании нужно задавать значение типа с момента создания стора.
В итоге решил воспользоваться приведением типов и выплевывать ошибки если что-то не так (type guards).
Насколько это правильно - не знаю.
*/

func (storage *MemStorage) Set(metricType string, name string, value interface{}) error {

	fmt.Println(metricType, name, value, reflect.TypeOf(value))

	key := fmt.Sprintf(`%s_%s`, metricType, name)

	if metricType == consts.COUNTER {
		if reflect.TypeOf(value).String() != "int64" {
			return errors.New("неверный тип value")
		}
		if storage.data[key] != nil {
			/* ВОПРОС ПРОВЕРЯЮЩЕМУ:
			Если типизировать следующим образом: storage.data[key].(consts.Counter) + value.(consts.Counter),
			То получаем панику:  interface conversion: interface {} is int64, not consts.Counter
			Хотя type Counter - это int64. Странная штуковина, я хочу обозвать именованным типом, а программа ругается.
			Это ограничение языка или я что-то не так делаю?
			*/
			storage.data[key] = storage.data[key].(int64) + value.(int64)
		} else {
			storage.data[key] = value
		}
		return nil
	}

	if metricType == consts.GAUGE {
		if reflect.TypeOf(value).String() != "float64" {
			return errors.New("неверный тип value")
		}
		storage.data[key] = value
		return nil
	}

	return errors.New("неизвестный тип метрики")
}

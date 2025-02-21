package storage

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemStorage_Set_MetricTypesCheck(t *testing.T) {
	type args struct {
		metricType consts.Metric
		name       string
		Gauge      *float64
		Counter    *int64
	}

	var gaugeVal = 111.1
	var counterVal int64 = 111

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Передаем метрику gauge cо значением Gauge - значение записалось в стор",
			args: args{
				metricType: consts.GAUGE,
				name:       "name",
				Gauge:      &gaugeVal,
			},
		},
		{
			name: "Передаем метрику gauge cо значением Counter - значение не записалось в стор",
			args: args{
				metricType: consts.GAUGE,
				name:       "name",
				Counter:    &counterVal,
			},
		},
		{
			name: "Передаем метрику counter cо значением int64 - значение записалось в стор",
			args: args{
				metricType: consts.COUNTER,
				name:       "name",
				Counter:    &counterVal,
			},
		},
		{
			name: "Передаем метрику counter cо значением float64 - значение не записалось в стор",
			args: args{
				metricType: consts.COUNTER,
				name:       "name",
				Gauge:      &gaugeVal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemStorage{
				data: MemStorageData{},
			}

			if tt.args.metricType == consts.GAUGE {
				storage.SetGauge(tt.args.name, tt.args.Gauge)
			} else if tt.args.metricType == consts.COUNTER {
				storage.SetCounter(tt.args.name, tt.args.Counter)
			}

			value := storage.Get(tt.args.metricType, tt.args.name)
			if tt.args.metricType == consts.GAUGE {
				if tt.args.Gauge != nil {
					assert.Equal(t, *tt.args.Gauge, value)
				}
				if tt.args.Gauge == nil {
					assert.Equal(t, nil, value)
				}
			}
			if tt.args.metricType == consts.COUNTER {
				if tt.args.Counter != nil {
					assert.Equal(t, *tt.args.Counter, value)
				}
				if tt.args.Counter == nil {
					assert.Equal(t, nil, value)
				}
			}
		})
	}
}

func TestMemStorage_Set_SaveGaugeMetricTwice(t *testing.T) {
	t.Run("Передаем метрику gauge cо значением float64 2 раза - записывается последнее переданное значение и имеет тип float64", func(t *testing.T) {
		storage := &MemStorage{
			data: MemStorageData{},
		}

		name := "test metric"
		gauge1 := 111.1
		gauge2 := 105.1
		expectedGauge := 105.1

		storage.SetGauge(name, &gauge1)
		storage.SetGauge(name, &gauge2)
		result := storage.Get(consts.GAUGE, name)

		assert.Equal(t, expectedGauge, result)
	})
}

func TestMemStorage_Set_SaveCounterMetricTwice(t *testing.T) {
	t.Run("Передаем метрику counter cо значением int64 2 раза - записывается сумма значений и имеет тип int64", func(t *testing.T) {
		storage := &MemStorage{
			data: MemStorageData{},
		}

		name := "test metric"
		var counter1 int64 = 100
		var counter2 int64 = 200
		var expectedCounter int64 = 300

		storage.SetCounter(name, &counter1)
		storage.SetCounter(name, &counter2)
		result := storage.Get(consts.COUNTER, name)

		assert.Equal(t, expectedCounter, result)
	})
}

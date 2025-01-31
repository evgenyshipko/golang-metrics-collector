package storage

import (
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/consts"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemStorage_Set_MetricTypesCheck(t *testing.T) {
	type args struct {
		metricType consts.Metric
		name       string
		value      interface{}
	}
	tests := []struct {
		name                 string
		args                 args
		wantErr              bool
		expectedErrorMessage string
	}{
		{
			name:    "Передаем метрику gauge cо значением float64 - ошибки нет",
			wantErr: false,
			args: args{
				metricType: consts.GAUGE,
				name:       "name",
				value:      111.1,
			},
		},
		{
			name:    "Передаем метрику gauge cо значением int - ошибки нет",
			wantErr: false,
			args: args{
				metricType: consts.GAUGE,
				name:       "name",
				value:      111,
			},
		},
		{
			name:    "Передаем метрику counter cо значением float64 - ошибка",
			wantErr: true,
			args: args{
				metricType: consts.COUNTER,
				name:       "name",
				value:      111.1,
			},
		},
		{
			name:    "Передаем метрику counter cо значением int - ошибки нет",
			wantErr: false,
			args: args{
				metricType: consts.COUNTER,
				name:       "name",
				value:      111,
			},
		},
		{
			name:    "Передаем метрику counter cо значением int64 - ошибки нет",
			wantErr: false,
			args: args{
				metricType: consts.COUNTER,
				name:       "name",
				value:      9223372036854775807,
			},
		},
		{
			name:                 "Передаем неизвестную метрику - ошибка",
			wantErr:              true,
			expectedErrorMessage: "неизвестный тип метрики",
			args: args{
				metricType: "unknown metric",
				name:       "name",
				value:      9223372036854775807,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemStorage{
				data: make(map[string]interface{}),
			}
			err := storage.Set(tt.args.metricType, tt.args.name, tt.args.value)
			fmt.Println("err", err, "tt.wantErr", tt.wantErr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.expectedErrorMessage != "" {
				assert.Equal(t, err.Error(), tt.expectedErrorMessage)
			}
		})
	}
}

func TestMemStorage_Set_SaveMetricTwice(t *testing.T) {
	type args struct {
		metricType    consts.Metric
		name          string
		value         interface{}
		expectedValue interface{}
	}
	tests := []struct {
		name                 string
		args                 args
		wantErr              bool
		expectedErrorMessage string
	}{
		{
			name:    "Передаем метрику gauge cо значением float64 2 раза - значение остается таким же и имеет тип float64",
			wantErr: false,
			args: args{
				metricType:    consts.GAUGE,
				name:          "name",
				value:         111.1,
				expectedValue: 111.1,
			},
		},
		{
			//FIXME: тест не помог найти ошибку в коде, надо переделать чтобы можно было 2 разные значения посылать
			name:    "Передаем метрику counter cо значением int 2 раза - значение складывается с предыдущим и имеет тип int64",
			wantErr: false,
			args: args{
				metricType:    consts.COUNTER,
				name:          "name",
				value:         111,
				expectedValue: int64(222),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemStorage{
				data: make(map[string]interface{}),
			}
			err := storage.Set(tt.args.metricType, tt.args.name, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
			err = storage.Set(tt.args.metricType, tt.args.name, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.expectedErrorMessage != "" {
				assert.Equal(t, err.Error(), tt.expectedErrorMessage)
			}

			key := fmt.Sprintf("%s_%s", tt.args.metricType, tt.args.name)
			assert.Equal(t, tt.args.expectedValue, storage.data[key])
		})
	}
}

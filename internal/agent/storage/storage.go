package storage

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/types"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"sync"
)

type MetricData struct {
	Value interface{}
	Type  consts.Metric
}

type MetricStorageData map[string]MetricData

type MetricStorage struct {
	Data MetricStorageData
	mu   *sync.RWMutex
}

func NewMetricStorage() *MetricStorage {
	storage := &MetricStorage{
		Data: make(MetricStorageData),
		mu:   &sync.RWMutex{},
	}
	return storage
}

func (s *MetricStorage) Set(data types.Data) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Data[data.Name] = MetricData{
		Value: data.Value,
		Type:  data.Type,
	}
}

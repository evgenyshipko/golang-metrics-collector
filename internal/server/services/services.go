package services

import (
	c "github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
	"time"
)

type Service interface {
	ProcessMetric(metricData c.MetricData) (c.MetricData, error)
	GetMetricData(metricData c.MetricData) (c.MetricData, int, error)
	GetMetricValue(metricData c.MetricData) (c.Values, int, error)
}

type MetricService struct {
	store           *storage.MemStorage
	storeInterval   time.Duration
	fileStoragePath string
}

func NewMetricService(store *storage.MemStorage, storeInterval time.Duration, fileStoragePath string) Service {
	return &MetricService{
		store:           store,
		storeInterval:   storeInterval,
		fileStoragePath: fileStoragePath,
	}
}

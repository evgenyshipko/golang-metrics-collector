package services

import (
	"context"
	"time"

	c "github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
)

type Service interface {
	ProcessMetric(ctx context.Context, metricData c.MetricData) (c.MetricData, error)
	GetMetricData(ctx context.Context, metricData c.MetricData) (c.MetricData, int, error)
	GetMetricValue(ctx context.Context, metricData c.MetricData) (c.Values, int, error)
	ProcessMetricArr(ctx context.Context, metricData []c.MetricData) error
}

type MetricService struct {
	store           storage.Storage
	storeInterval   time.Duration
	fileStoragePath string
}

func NewMetricService(store storage.Storage, storeInterval time.Duration, fileStoragePath string) Service {
	return &MetricService{
		store:           store,
		storeInterval:   storeInterval,
		fileStoragePath: fileStoragePath,
	}
}

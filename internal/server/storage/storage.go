package storage

import (
	"context"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/db"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/setup"
)

type Data struct {
	consts.Values
	Name string `json:"name"`
}

type StorageData []Data

type Storage interface {
	Get(ctx context.Context, metricType consts.Metric, name string) *consts.Values
	GetAll(ctx context.Context) (*StorageData, error)
	SetGauge(ctx context.Context, name string, value *float64)
	SetCounter(ctx context.Context, name string, value *int64)
	SetData(ctx context.Context, data StorageData) error
	IsAvailable(ctx context.Context) bool
	Close() error
}

func NewStorage(cfg *setup.ServerStartupValues) (Storage, error) {
	if cfg.DatabaseDSN != "" {
		conn, err := db.ConnectToDB(cfg.DatabaseDSN, cfg.AutoMigrations)
		if err != nil {
			return &MemStorage{}, err
		}
		return NewSQLStorage(conn, cfg), nil
	} else {
		return NewMemStorage(), nil
	}
}

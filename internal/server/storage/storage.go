package storage

import (
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
	Get(metricType consts.Metric, name string) *consts.Values
	GetAll() (*StorageData, error)
	SetGauge(name string, value *float64)
	SetCounter(name string, value *int64)
	SetData(data StorageData) error
	IsAvailable() bool
	Close() error
}

func NewStorage(cfg *setup.ServerStartupValues) (Storage, error) {
	if cfg.DatabaseDSN != "" {
		conn, err := db.ConnectToDB(cfg.DatabaseDSN)
		if err != nil {
			return &MemStorage{}, err
		}
		return NewSQLStorage(conn, cfg), nil
	} else {
		return NewMemStorage(), nil
	}
}

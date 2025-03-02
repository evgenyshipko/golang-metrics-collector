package storage

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/db"
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

func NewStorage(databaseDSN string) (Storage, error) {
	if databaseDSN != "" {
		conn, err := db.ConnectToDB(databaseDSN)
		if err != nil {
			return &MemStorage{}, err
		}
		return NewSQLStorage(conn), nil
	} else {
		return NewMemStorage(), nil
	}
}

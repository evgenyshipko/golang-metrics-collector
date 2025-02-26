package db

import (
	"database/sql"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	_ "github.com/jackc/pgx/stdlib"
)

func ConnectToDb(serverDSN string) *sql.DB {
	db, err := sql.Open("pgx", serverDSN)
	if err != nil {
		logger.Instance.Warnw("ConnectToDb", "Не удалось подключиться к базе данных", err)
	}
	return db
}

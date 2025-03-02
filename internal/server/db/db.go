package db

import (
	"database/sql"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/jackc/pgx/stdlib"
)

func ConnectToDB(serverDSN string) (*sql.DB, error) {
	db, err := sql.Open("pgx", serverDSN)
	if err != nil {
		logger.Instance.Warnw("ConnectToDB", "Не удалось подключиться к базе данных", err)
		return nil, err
	}
	return db, nil
}

type SQLExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
}

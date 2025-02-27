package db

import (
	"database/sql"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/jackc/pgx/stdlib"
)

func ConnectToDB(serverDSN string) *sql.DB {
	db, err := sql.Open("pgx", serverDSN)
	if err != nil {
		logger.Instance.Warnw("ConnectToDB", "Не удалось подключиться к базе данных", err)
	}

	err = RunMigrations(db)
	if err != nil {
		logger.Instance.Warnw("RunMigrations", "Ошибка проката миграций", err)
	}

	return db
}

// TODO: по-хорошему не код приложения должен запускать матгации, а отдельная джоба. Запуск вызовет проблемы если сделать множество экземпляров сервера.
func RunMigrations(db *sql.DB) error {

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/server/db/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	logger.Instance.Info("Migrations applied successfully!")
	return nil
}

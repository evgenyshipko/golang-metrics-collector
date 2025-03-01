package retry

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"net"
	"os"
	"time"
)

func WithRetry[T any](fn func() (T, error)) (T, error) {
	var result T
	retryIntervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	var err error
	for i, wait := range retryIntervals {
		result, err = fn()
		if err == nil || !isRetriableError(err) {
			logger.Instance.Infow("Retry срабатывает без ошибки")
			return result, err
		}
		logger.Instance.Warnw("WithRetry", "ошибка", err)

		logger.Instance.Warnf("Попытка %d, ждем %s перед следующим запросом...\n", i+1, wait)
		time.Sleep(wait)
	}
	return result, err
}

func ExecuteTransactionWithRetry(db *sql.DB, fn func(tx *sql.Tx) error) error {
	retryIntervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	for attempt, interval := range retryIntervals {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // ⏳ Тайм-аут 3 секунды

		defer cancel()

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			logger.Instance.Warnw("ExecuteTransactionWithRetry", "ошибка создания транзакции", err)
			if isRetriableError(err) {
				logger.Instance.Warnf("ExecuteTransactionWithRetry Попытка %d не удалась, ждем %s перед повтором...", attempt+1, interval)
				time.Sleep(interval)
				continue
			}
			return err
		}

		err = fn(tx)
		if err != nil {
			logger.Instance.Warnw("ExecuteTransactionWithRetry", "ошибка функции", err)
			tx.Rollback()
			if isRetriableError(err) {
				logger.Instance.Warnf("ExecuteTransactionWithRetry Попытка %d не удалась, ждем %s перед повтором...", attempt+1, interval)
				time.Sleep(interval)
				continue
			}
			return err
		}

		return tx.Commit()
	}

	return fmt.Errorf("все попытки завершились неудачей")
}

func isRetriableError(err error) bool {
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		logger.Instance.Infow("isRetriableError", "Ошибка сети", err)
		return true
	}

	var syscallErr *os.SyscallError
	if errors.As(err, &syscallErr) {
		logger.Instance.Infow("isRetriableError", "Ошибка системного вызова", err)
		return true
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.DeadlockDetected, // deadlock detected
			pgerrcode.SerializationFailure, // could not serialize access
			pgerrcode.LockNotAvailable,     // lock is not available
			pgerrcode.ConnectionException,  // connection issues
			pgerrcode.TooManyConnections,   // too many connections
			pgerrcode.AdminShutdown,        // сервер перезапущен админом
			pgerrcode.CrashShutdown,        // сервер упал и перезапустился
			pgerrcode.IOError,              // ошибка ввода-вывода
			pgerrcode.QueryCanceled,        // клиент прервал запрос
			pgerrcode.ConnectionFailure:
			return true
		default:
			return false
		}
	}
	return false
}

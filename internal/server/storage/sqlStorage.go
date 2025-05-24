package storage

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/retry"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/setup"
)

type SQLStorage struct {
	cfg        *setup.ServerStartupValues
	db         *sql.DB
	statements map[string]*sql.Stmt
	mu         sync.RWMutex //ЗАПОМНИТЬ: мапа не потокобезопасна, поэтому при конкуррентном чтении/записи могут возникать ошибки конкуррентного доступа к данным
}

func NewSQLStorage(db *sql.DB, cfg *setup.ServerStartupValues) *SQLStorage {
	return &SQLStorage{
		cfg:        cfg,
		db:         db,
		statements: map[string]*sql.Stmt{},
		mu:         sync.RWMutex{},
	}
}

func (storage *SQLStorage) prepareStmt(query string) (*sql.Stmt, error) {
	storage.mu.RLock() // Блокируем только для чтения
	stmt, exists := storage.statements[query]
	storage.mu.RUnlock() // Разблокируем чтение
	if exists {
		return stmt, nil
	}

	storage.mu.Lock()         // Блокируем на запись
	defer storage.mu.Unlock() // Разблокируем запись

	stmt, err := storage.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	storage.statements[query] = stmt
	return stmt, nil
}

func (storage *SQLStorage) Get(ctx context.Context, metricType consts.Metric, name string) *consts.Values {
	values, err := retry.WithRetry(func() (consts.Values, error) {

		query := "SELECT value_int, value_float FROM metrics WHERE name = $1 AND type = $2"

		stmt, err := storage.prepareStmt(query)
		if err != nil {
			return consts.Values{}, err
		}

		row := stmt.QueryRowContext(ctx, name, metricType)

		var values consts.Values

		err = row.Scan(&values.Counter, &values.Gauge)
		return values, err
	}, storage.cfg.RetryIntervals)

	if err != nil {
		logger.Instance.Warnw("NewSQLStorage", "Get", err)
		return &consts.Values{}
	}

	return &values
}

func (storage *SQLStorage) SetGauge(ctx context.Context, name string, value *float64) {
	_, err := retry.WithRetry(func() (string, error) {
		innerErr := storage.insertData(ctx, nil, name, consts.GAUGE, value, nil)
		return "", innerErr
	}, storage.cfg.RetryIntervals)
	if err != nil {
		logger.Instance.Warnw("NewSQLStorage", "SetGauge", err)
	}
}

func (storage *SQLStorage) SetCounter(ctx context.Context, name string, value *int64) {
	_, err := retry.WithRetry(func() (string, error) {
		innerErr := storage.insertData(ctx, nil, name, consts.COUNTER, nil, value)
		return "", innerErr
	}, storage.cfg.RetryIntervals)
	if err != nil {
		logger.Instance.Warnw("NewSQLStorage", "SetCounter", err)
	}
}

func (storage *SQLStorage) insertData(ctx context.Context, tx *sql.Tx, name string, metricType consts.Metric, valueFloatPointer *float64, valueIntPointer *int64) error {

	logger.Instance.Debugw("insertData", "tx", tx, "name", name, "metricType", metricType, "valueFloat", valueFloatPointer, "valueInt", valueIntPointer)

	query := `
    INSERT INTO metrics (name, type, value_int, value_float)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (name, type) DO UPDATE 
    SET 
        value_int = CASE 
            WHEN EXCLUDED.type = 'counter' THEN metrics.value_int + EXCLUDED.value_int 
            ELSE metrics.value_int 
        END,
        value_float = CASE 
            WHEN EXCLUDED.type = 'gauge' THEN EXCLUDED.value_float 
            ELSE metrics.value_float 
        END;
`
	logger.Instance.Debug(debugQuery(query, name, metricType, valueIntPointer, valueFloatPointer))

	if tx != nil {
		// если запрос в рамках транзакции, то тогда запрос не подготавливаем
		_, err := tx.ExecContext(ctx, query, name, metricType, valueIntPointer, valueFloatPointer)
		if err != nil {
			return err
		}
		return nil
	}

	stmt, err := storage.prepareStmt(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, name, metricType, valueIntPointer, valueFloatPointer)
	if err != nil {
		return err
	}
	return nil
}

func (storage *SQLStorage) SetData(ctx context.Context, data StorageData) error {

	err := retry.ExecuteTransactionWithRetry(ctx, storage.db, func(tx *sql.Tx) error {
		for _, val := range data {
			var metricType consts.Metric
			if val.Counter != nil {
				metricType = consts.COUNTER
			} else if val.Gauge != nil {
				metricType = consts.GAUGE
			}

			err := storage.insertData(ctx, tx, val.Name, metricType, val.Gauge, val.Counter)
			if err != nil {
				return err
			}
		}
		return nil
	}, storage.cfg.RetryIntervals)

	if err != nil {
		return err
	}

	return nil
}

func (storage *SQLStorage) GetAll(ctx context.Context) (*StorageData, error) {

	metrics := StorageData{}

	rows, err := retry.WithRetry(func() (*sql.Rows, error) {

		query := "SELECT name, value_int, value_float from metrics"

		stmt, err := storage.prepareStmt(query)
		if err != nil {
			return nil, err
		}

		return stmt.QueryContext(ctx)
	}, storage.cfg.RetryIntervals)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var v Data
		err = rows.Scan(&v.Name, &v.Counter, &v.Gauge)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, v)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &metrics, nil
}

func (storage *SQLStorage) IsAvailable(ctx context.Context) bool {
	err := storage.db.PingContext(ctx)
	if err != nil {
		logger.Instance.Warnw("IsAvailable", "err", err)
		return false
	}
	return true
}

func (storage *SQLStorage) Close() error {
	return storage.db.Close()
}

func debugQuery(query string, args ...interface{}) string {
	// Регулярка для поиска плейсхолдеров $1, $2, ...
	re := regexp.MustCompile(`\$\d+`)

	// Индекс текущего аргумента
	argIndex := 0

	// Заменяем $1, $2, ... на реальные значения
	result := re.ReplaceAllStringFunc(query, func(_ string) string {
		if argIndex >= len(args) {
			return "NULL" // Если аргументов меньше, чем плейсхолдеров
		}
		arg := args[argIndex]
		argIndex++

		// Форматируем аргументы в SQL-friendly строку
		switch v := arg.(type) {
		case string:
			return "'" + strings.ReplaceAll(v, "'", "''") + "'" // SQL-экранирование кавычек
		case []byte:
			return "'" + strings.ReplaceAll(string(v), "'", "''") + "'"
		case int, int64, float64:
			return fmt.Sprintf("%v", v)
		case bool:
			return strconv.FormatBool(v)
		case nil:
			return "NULL"
		default:
			return fmt.Sprintf("'%v'", v) // Для остальных типов
		}
	})

	return result
}

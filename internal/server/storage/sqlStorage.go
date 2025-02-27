package storage

import (
	"database/sql"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"regexp"
	"strconv"
	"strings"
)

type SQLStorage struct {
	db *sql.DB
}

func NewSQLStorage(db *sql.DB) *SQLStorage {
	return &SQLStorage{
		db: db,
	}
}

func (storage *SQLStorage) Get(metricType consts.Metric, name string) *consts.Values {

	row := storage.db.QueryRow("SELECT value_int, value_float FROM metrics WHERE name = $1 AND type = $2", name, metricType)

	var values consts.Values

	err := row.Scan(&values.Counter, &values.Gauge)
	if err != nil {
		logger.Instance.Warnw("NewSQLStorage", "Get", err)
		return &consts.Values{}
	}

	return &values
}

func (storage *SQLStorage) SetGauge(name string, value *float64) {
	err := storage.insertData(name, consts.GAUGE, value, nil)
	if err != nil {
		logger.Instance.Warnw("NewSQLStorage", "SetGauge", err)
	}
}

func (storage *SQLStorage) SetCounter(name string, value *int64) {
	err := storage.insertData(name, consts.COUNTER, nil, value)
	if err != nil {
		logger.Instance.Warnw("NewSQLStorage", "SetCounter", err)
	}
}

func (storage *SQLStorage) insertData(name string, metricType consts.Metric, valueFloatPointer *float64, valueIntPointer *int64) error {

	query := `
    INSERT INTO metrics (name, type, value_int, value_float)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (name, type) DO UPDATE 
    SET 
        value_int = CASE 
            WHEN EXCLUDED.type = 'counter' THEN EXCLUDED.value_int 
            ELSE metrics.value_int 
        END,
        value_float = CASE 
            WHEN EXCLUDED.type = 'gauge' THEN EXCLUDED.value_float 
            ELSE metrics.value_float 
        END;
`
	logger.Instance.Debug(debugQuery(query, name, metricType, valueIntPointer, valueFloatPointer))

	_, err := storage.db.Exec(query, name, metricType, valueIntPointer, valueFloatPointer)

	if err != nil {
		return err
	}
	return nil
}

func (storage *SQLStorage) SetData(data StorageData) {
	for _, value := range data {
		if value.Counter != nil {
			storage.SetCounter(value.Name, value.Counter)
		}
		if value.Gauge != nil {
			storage.SetGauge(value.Name, value.Gauge)
		}
	}
}

func (storage *SQLStorage) GetAll() (*StorageData, error) {

	metrics := StorageData{}

	rows, err := storage.db.Query("SELECT name, value_int, value_float from metrics")
	if err != nil {
		return nil, err
	}

	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var v Data
		err = rows.Scan(&v.Name, &v.Counter, &v.Gauge)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, v)
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &metrics, nil
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

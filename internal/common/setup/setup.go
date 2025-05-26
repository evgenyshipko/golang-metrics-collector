package setup

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
)

func GetInterval(envName string, flagVal *int, validate ...bool) (time.Duration, error) {
	validateZero := true
	if len(validate) > 0 {
		validateZero = validate[0]
	}

	envInterval, exists := os.LookupEnv(envName)
	intInterval := 0
	if exists {
		val, err := strconv.Atoi(envInterval)
		if err != nil {
			logger.Instance.Warn(fmt.Sprintf("ошибка конвертации енва %s, будем брать из флагов", envName))
		}
		intInterval = val
	} else {
		intInterval = *flagVal
	}
	err := validateInterval(intInterval, validateZero)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	return intToSeconds(intInterval), nil
}

func intToSeconds(num int) time.Duration {
	return time.Duration(num) * time.Second
}

func validateInterval(num int, validateZero bool) error {
	if validateZero && num <= 0 {
		return errors.New("интервал должен быть положительным и больше нуля")
	}
	if !validateZero && num < 0 {
		return errors.New("интервал должен быть положительным")
	}
	return nil
}

func GetStringVariable(envName string, flagVal *string) string {
	env, exists := os.LookupEnv(envName)

	if exists {
		return env
	}
	return *flagVal
}

func GetBoolVariable(envName string, flagVal *bool) bool {
	env, exists := os.LookupEnv(envName)
	if exists {
		boolean, err1 := strconv.ParseBool(env)
		if err1 != nil {
			return *flagVal
		}
		return boolean
	}
	return *flagVal
}

func GetIntervals(envName string, flagVal *string) ([]time.Duration, error) {
	env, exists := os.LookupEnv(envName)
	if exists {
		return convertStringToDurations(env)
	}
	return convertStringToDurations(*flagVal)
}

func GetIntVariable(envName string, flagVal *int) (int, error) {
	env, exists := os.LookupEnv(envName)
	if exists {
		int64_, err := strconv.ParseInt(env, 10, 32)
		if err != nil {
			return 0, err
		}
		return int(int64_), nil
	}
	return *flagVal, nil
}

func convertStringToDurations(input string) ([]time.Duration, error) {
	intervals := strings.Split(input, ",")

	var durations []time.Duration

	for _, interval := range intervals {
		interval = strings.TrimSpace(interval)

		duration, err := time.ParseDuration(interval)
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования '%s' в time.Duration: %v", interval, err)
		}

		durations = append(durations, duration)
	}

	return durations, nil
}

package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/tasks"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/types"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	dataCh := make(chan types.MetricMessage, 100)
	errCh := make(chan error)

	defer func() {
		close(dataCh)
		close(errCh)
		close(signalChan)
		logger.Sync()
	}()

	vars, err := setup.GetStartupValues(os.Args[1:])
	if err != nil {
		logger.Instance.Fatalw("Аргументы не прошли валидацию", err)
	}

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go tasks.SendMetricsTask(vars, dataCh, errCh)

	go tasks.MetricsGenerator(vars.PollInterval, dataCh)

	go tasks.AdditionalMetricsGenerator(vars.PollInterval, dataCh)

	go tasks.ErrorsConsumer(errCh)

	go tasks.LogChanLength(dataCh)

	// Ожидаем сигнала завершения
	<-signalChan

	// Даём время горутине завершиться
	time.Sleep(1 * time.Second)
	logger.Instance.Info("Агент завершил работу")
}

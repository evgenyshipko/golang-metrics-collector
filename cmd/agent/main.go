package main

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/tasks"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/types"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	vars, err := setup.GetStartupValues(os.Args[1:])
	if err != nil {
		logger.Instance.Fatalw("Аргументы не прошли валидацию", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	dataCh := make(chan types.ChanData, 100)

	errCh := make(chan error)

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

	defer func() {
		close(dataCh)
		close(errCh)
		close(signalChan)
		logger.Sync()
	}()
}

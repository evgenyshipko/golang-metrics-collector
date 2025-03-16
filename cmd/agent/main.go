package main

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/storage"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/tasks"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/types"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	defer logger.Sync()

	vars, err := setup.GetStartupValues(os.Args[1:])
	if err != nil {
		logger.Instance.Fatalw("Аргументы не прошли валидацию", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	metricStorage := storage.NewMetricStorage()

	go tasks.SendMetricsTask(vars, metricStorage)

	dataCh := make(chan types.ChanData)

	defer close(dataCh)

	go tasks.MetricsGenerator(vars.PollInterval, dataCh)

	go tasks.AdditionalMetricsGenerator(vars.PollInterval, dataCh)

	go tasks.MetricsConsumer(metricStorage, dataCh)

	// Ожидаем сигнала завершения
	<-signalChan

	// Даём время горутине завершиться
	time.Sleep(1 * time.Second)
	logger.Instance.Info("Агент завершил работу")
}

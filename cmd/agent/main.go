package main

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/storage"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/tasks"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	log := logger.InitLogger()

	defer log.Sync()

	vars, err := setup.GetStartupValues()
	if err != nil {
		logger.Fatal("Аргументы не прошли валидацию", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	metrics := storage.NewMetricStorage()

	go tasks.CollectMetricsTask(vars.PollInterval, &metrics)

	go tasks.SendMetricsTask(vars.ReportInterval, &metrics, vars.Host)

	// Ожидаем сигнала завершения
	<-signalChan

	// Даём время горутине завершиться
	time.Sleep(1 * time.Second)
	logger.Info("Агент завершил работу")
}

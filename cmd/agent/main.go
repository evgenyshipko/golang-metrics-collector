package main

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/tasks"
	"github.com/evgenyshipko/golang-metrics-collector/internal/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	vars := setup.GetStartupValues()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	metrics := map[string]interface{}{}

	go tasks.CollectMetricsTask(time.Duration(vars.PollInterval)*time.Second, &metrics)

	go tasks.SendMetricsTask(time.Duration(vars.ReportInterval)*time.Second, &metrics, vars.Host)

	// Ожидаем сигнала завершения
	<-signalChan

	// Даём время горутине завершиться
	time.Sleep(1 * time.Second)
	logger.Info("Агент завершил работу")
}

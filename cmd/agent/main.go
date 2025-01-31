package main

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/tasks"
	"github.com/evgenyshipko/golang-metrics-collector/internal/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	metrics := map[string]interface{}{}

	go tasks.CollectMetricsTask(2*time.Second, &metrics)

	go tasks.SendMetricsTask(10*time.Second, &metrics)

	// Ожидаем сигнала завершения
	<-signalChan

	// Даём время горутине завершиться
	time.Sleep(1 * time.Second)
	logger.Info("Агент завершил работу")
}

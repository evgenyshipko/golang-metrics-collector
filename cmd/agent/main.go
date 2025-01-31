package main

import (
	"flag"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/tasks"
	"github.com/evgenyshipko/golang-metrics-collector/internal/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	metricsHost := flag.String("a", "localhost:8080", "metric server host")

	reportInterval := flag.Int("r", 10, "interval between report metrics")

	pollInterval := flag.Int("p", 2, "interval between polling metrics")

	flag.Parse()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	metrics := map[string]interface{}{}

	go tasks.CollectMetricsTask(time.Duration(*pollInterval)*time.Second, &metrics)

	go tasks.SendMetricsTask(time.Duration(*reportInterval)*time.Second, &metrics, *metricsHost)

	// Ожидаем сигнала завершения
	<-signalChan

	// Даём время горутине завершиться
	time.Sleep(1 * time.Second)
	logger.Info("Агент завершил работу")
}

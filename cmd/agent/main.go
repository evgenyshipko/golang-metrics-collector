package main

import (
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/tasks"
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
	fmt.Println("Агент завершил работу")
}

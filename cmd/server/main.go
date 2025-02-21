package main

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/server"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/tasks"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer logger.Sync()

	values, err := setup.GetStartupValues(os.Args[1:])
	if err != nil {
		logger.Instance.Fatalw("Аргументы не прошли валидацию", err)
	}

	customServer := server.Create(&values)

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)

	go customServer.Start()

	go tasks.WriteMetricsToFileTask(values.StoreInterval, values.FileStoragePath, customServer.GetStoreData())

	<-stopSignal

	customServer.ShutDown()
}

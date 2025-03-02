package main

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/files"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/server"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/tasks"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	values, err := setup.GetStartupValues(os.Args[1:])
	if err != nil {
		logger.Instance.Fatalw("Аргументы не прошли валидацию", err)
	}

	store, err := storage.NewStorage(&values)
	if err != nil {
		logger.Instance.Warnw("server.Create", "ошибка создания store", err)
		return
	}

	if values.Restore {
		files.RecoverFromFile(values.FileStoragePath, store, values.RetryIntervals)
	}

	customServer := server.Create(&values, store)

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)

	go customServer.Start()

	go tasks.WriteMetricsToFileTask(values.StoreInterval, values.FileStoragePath, customServer)

	<-stopSignal

	customServer.ShutDown()

	defer func() {
		logger.Sync()
		store.Close()
	}()
}

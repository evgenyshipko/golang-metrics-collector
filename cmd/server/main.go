package main

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/files"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/grpcServer"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/services"
	"os"
	"os/signal"
	"syscall"

	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/server"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/tasks"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {

	logger.Instance.Infof("\nBuild version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

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

	metricService := services.NewMetricService(store, values.StoreInterval, values.FileStoragePath)

	customServer := server.Create(&values, store, metricService)

	// через этот канал сообщим основному потоку, что соединения закрыты
	idleConnsClosed := make(chan struct{})
	// канал для перенаправления прерываний
	// поскольку нужно отловить всего одно прерывание,
	// ёмкости 1 для канала будет достаточно
	stopSignal := make(chan os.Signal, 1)
	// регистрируем перенаправление прерываний
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// РАЗОБРАТЬСЯ: сделал по примеру из урока. Нужно ли обязательно делать горутину, как сделано ниже или и без нее было норм?
	go func() {
		// читаем из канала прерываний
		// поскольку нужно прочитать только одно прерывание,
		// можно обойтись без цикла
		<-stopSignal
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		customServer.ShutDown()

		// сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()

	go customServer.Start()

	go grpcServer.StartGrpcServer(metricService)

	go tasks.WriteMetricsToFileTask(values.StoreInterval, values.FileStoragePath, customServer)

	// ждём завершения процедуры graceful shutdown
	<-idleConnsClosed

	defer func() {
		logger.Sync()
		store.Close()
	}()
}

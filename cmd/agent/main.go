package main

import (
	"context"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/tasks"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/types"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"os"
	"os/signal"
	"syscall"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {

	logger.Instance.Infof("\nBuild version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	dataCh := make(chan types.MetricMessage, 100)
	errCh := make(chan error)

	defer func() {
		stop()
		close(dataCh)
		close(errCh)
		logger.Sync()
	}()

	vars, err := setup.GetStartupValues(os.Args[1:])
	if err != nil {
		logger.Instance.Fatalw("Аргументы не прошли валидацию", err)
	}

	go tasks.SendMetricsTask(ctx, vars, dataCh, errCh)

	go tasks.MetricsGenerator(ctx, vars.PollInterval, dataCh)

	go tasks.AdditionalMetricsGenerator(ctx, vars.PollInterval, dataCh)

	go tasks.ErrorsConsumer(errCh)

	go tasks.LogChanLength(dataCh)

	// Ожидаем сигнала завершения
	<-ctx.Done()

	logger.Instance.Info("Агент завершил работу")
}

package tasks

import (
	"context"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/converter"
	"sync"
	"time"

	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/requests"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/types"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
)

func SendMetricsTask(ctx context.Context, cfg setup.AgentStartupValues, dataChan <-chan types.MetricMessage, errChan chan<- error) {
	ticker := time.NewTicker(cfg.ReportInterval)

	requester := requests.NewRequester(cfg)

	var wg sync.WaitGroup

	for w := 1; w <= cfg.RateLimit; w++ {
		wg.Add(1)
		go worker(ctx, w, dataChan, errChan, requester, ticker, &wg)
	}

	go func() {
		<-ctx.Done()
		logger.Instance.Debug("SendMetricsTask <-ctx.Done()")
		ticker.Stop() // останавливаем тикер
		requester.Close()
	}()

	wg.Wait() // Ждём завершения всех горутин
}

/*
worker обрабатывает метрики с graceful shutdown при получении сигнала завершения.
Все блокирующие операции защищены select с проверкой ctx.Done().
*/
func worker(ctx context.Context, id int, jobs <-chan types.MetricMessage, errChan chan<- error, requester requests.Requester,
	ticker *time.Ticker, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %d starting\n", id)
	for {
		select {
		case <-ctx.Done(): // Получен сигнал завершения
			logger.Instance.Debugf("Worker %d exiting 1\n", id)
			return
		case <-ticker.C:
			for {
				select {
				case job, ok := <-jobs:
					if !ok { // Канал закрыт
						logger.Instance.Debugf("jobs channel closed\n")
						return
					}
					if job.Err != nil {
						errChan <- job.Err
						logger.Instance.Warnw("Обработка ошибки", "error", job.Err)
						continue
					}

					requestData, err := converter.GenerateMetricData(job.Data.Type, job.Data.Name, job.Data.Value)
					if err != nil {
						select {
						case errChan <- err:
						case <-ctx.Done():
							logger.Instance.Debugf("Worker %d exiting 2\n", id)
							return
						}
					}

					err = requester.SendMetric(requestData)
					if err != nil {
						select {
						case errChan <- err:
						case <-ctx.Done():
							logger.Instance.Debugf("Worker %d exiting 3\n", id)
							return
						}
					}
				case <-ctx.Done():
					logger.Instance.Debugf("Worker %d exiting 4\n", id)
					return
				}
			}
		}
	}
}

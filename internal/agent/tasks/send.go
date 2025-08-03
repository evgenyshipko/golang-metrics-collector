package tasks

import (
	"context"
	"fmt"
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
	}()

	wg.Wait() // Ждём завершения всех горутин
}

/*
worker обрабатывает метрики с graceful shutdown при получении сигнала завершения.
Все блокирующие операции защищены select с проверкой ctx.Done().
*/
func worker(ctx context.Context, id int, jobs <-chan types.MetricMessage, errChan chan<- error, requester *requests.Requester,
	ticker *time.Ticker, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %d starting\n", id)
	for {
		select {
		case <-ctx.Done(): // Получен сигнал завершения
			return
		case <-ticker.C:
			select {
			case job, ok := <-jobs:
				if !ok { // Канал закрыт
					return
				}
				if job.Err != nil {
					logger.Instance.Warnw("Обработка ошибки", "error", job.Err)
					continue
				}

				err := requester.SendMetric(job.Data.Type, job.Data.Name, job.Data.Value)
				if err != nil {
					select {
					case errChan <- err:
					case <-ctx.Done():
						return
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

package tasks

import (
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/requests"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/types"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"sync"
	"time"
)

func SendMetricsTask(cfg setup.AgentStartupValues, dataChan <-chan types.MetricMessage, errChan chan<- error) {
	ticker := time.NewTicker(cfg.ReportInterval)

	requester := requests.NewRequester(cfg)

	var wg sync.WaitGroup

	for w := 1; w <= cfg.RateLimit; w++ {
		wg.Add(1)
		go worker(w, dataChan, errChan, requester, ticker, &wg)
	}

	go func() {
		wg.Wait()     // Ждём завершения всех горутин
		ticker.Stop() // останавливаем тикер
	}()
}

func worker(id int, jobs <-chan types.MetricMessage, errChan chan<- error, requester *requests.Requester, ticker *time.Ticker, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %d starting\n", id)
	for range ticker.C {
		fmt.Printf("Worker %d processing metrics...\n", id)
		for job := range jobs {
			if job.Err != nil {
				logger.Instance.Warnw("Обработка ошибки", "error", job.Err)
				continue
			}

			logger.Instance.Infow("worker", "рабочий", id, "запущена задача", job)

			err := requester.SendMetric(job.Data.Type, job.Data.Name, job.Data.Value)
			if err != nil {
				errChan <- err
			}

			logger.Instance.Infow("worker", "рабочий", id, "закончил задачy", job)
		}
	}
}

package tasks

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/requests"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/types"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"time"
)

func SendMetricsTask(cfg setup.AgentStartupValues, dataChan <-chan types.ChanData, errChan chan<- error) {
	ticker := time.NewTicker(cfg.ReportInterval)
	defer ticker.Stop()

	requester := requests.NewRequester(cfg)

	for range ticker.C {
		for w := 1; w <= cfg.RateLimit; w++ {
			go worker(w, dataChan, errChan, requester)
		}
	}
}

func worker(id int, jobs <-chan types.ChanData, errChan chan<- error, requester *requests.Requester) {
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

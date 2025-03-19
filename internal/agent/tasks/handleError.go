package tasks

import "github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"

func ErrorsConsumer(errChan <-chan error) {
	for err := range errChan {
		if err != nil {
			logger.Instance.Warnw("ErrorsConsumer", "err", err)
		}
	}
}

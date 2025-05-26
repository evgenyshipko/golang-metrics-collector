package tasks

import (
	"time"

	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/types"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
)

func LogChanLength(ch chan types.MetricMessage) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		logger.Instance.Infof("Current ch length: %d\n", len(ch))
	}
}

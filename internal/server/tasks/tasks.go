package tasks

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/files"
	"time"

	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/server"
)

func WriteMetricsToFileTask(interval time.Duration, filePath string, server *server.CustomServer) {
	if interval == 0 {
		return
	}

	data, err := server.GetStoreData()
	if err != nil {
		logger.Instance.Warnw("GetStoreData", "не удалось получить данные", err)
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		files.WriteToFile(filePath, data)
	}
}

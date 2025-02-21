package tasks

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/files"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
	"time"
)

func WriteMetricsToFileTask(interval time.Duration, filePath string, data *storage.MemStorageData) {
	if interval == 0 {
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		files.WriteToFile(filePath, data)
	}
}

package requests

import (
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/consts"
	"github.com/go-resty/resty/v2"
)

func SendMetric(metricType consts.Metric, name string, value string) error {
	domain := "localhost:8080"
	url := fmt.Sprintf("http://%s/update/%s/%s/%s", domain, metricType, name, value)

	client := resty.New()
	resp, err := client.R().Post(url)

	if resp.StatusCode() == 200 {
		logger.Info("Метрики успешно отправлены")
		return nil
	}

	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос: %w", err)
	}

	logger.Info("SendMetric response", "url", url, "status", resp.Status(), "body", resp.Body())

	return nil
}

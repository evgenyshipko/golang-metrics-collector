package requests

import (
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/go-resty/resty/v2"
)

func SendMetric(domain string, metricType consts.Metric, name string, value string) error {
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

	logger.Info("SendMetric Response", "url", url, "status", resp.Status(), "body", resp.Body())

	return nil
}

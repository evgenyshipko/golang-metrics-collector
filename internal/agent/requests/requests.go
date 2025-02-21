package requests

import (
	"encoding/json"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/gzip"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/converter"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/go-resty/resty/v2"
)

func SendMetric(domain string, metricType consts.Metric, name string, value interface{}) error {

	requestData, err := converter.GenerateMetricData(metricType, name, value)
	if err != nil {
		logger.Instance.Warnw("SendMetric", "GenerateMetricData", err)
		return err
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		logger.Instance.Warnw("SendMetric", "json.Marshal err", err)
		return err
	}

	compressedBody, err := gzip.Compress(body)
	if err != nil {
		logger.Instance.Warnw("SendMetric", "compress err", err)
		return err
	}

	url := fmt.Sprintf("http://%s/update/", domain)

	client := resty.New()

	//ЗАПОМНИТЬ: resty автоматически добавляет заголовок "Accept-Encoding", "gzip" и распаковывает ответ если он пришел в gzip
	resp, err := client.R().
		SetBody(compressedBody).
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json").
		Post(url)

	if resp.StatusCode() == 200 {
		logger.Instance.Info("Метрики успешно отправлены")
	}

	if err != nil {
		logger.Instance.Errorf("SendMetric", "не удалось выполнить запрос", err)
	}

	respBody := resp.Body()

	logger.Instance.Infow("SendMetric Response", "url", url, "status", resp.Status(), "body", respBody)

	return nil
}

package requests

import (
	"encoding/json"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/converter"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/go-resty/resty/v2"
)

func SendMetric(domain string, metricType consts.Metric, name string, value interface{}) error {

	requestData, err := converter.GenerateMetricData(metricType, name, value)
	if err != nil {
		logger.Instance.Warnw("GenerateMetricData", err)
		return err
	}

	logger.Instance.Debug("SendMetric", "requestData", requestData)

	body, err := json.Marshal(requestData)
	if err != nil {
		logger.Instance.Warnw("SendMetric json.Marshal err", err)
		return err
	}

	url := fmt.Sprintf("http://%s/update/", domain)

	client := resty.New()
	resp, err := client.R().SetBody(body).Post(url)

	if resp.StatusCode() == 200 {
		logger.Instance.Info("Метрики успешно отправлены")
		return nil
	}

	if err != nil {
		logger.Instance.Errorf("не удалось выполнить запрос: \n%w", err)
	}

	logger.Instance.Infow("SendMetric Response", "url", url, "status", resp.Status(), "body", resp.Body())

	return nil
}

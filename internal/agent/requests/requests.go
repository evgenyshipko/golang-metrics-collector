package requests

import (
	"encoding/json"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/converter"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/gzip"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"time"
)

var retryAfterFunc resty.RetryAfterFunc = func(c *resty.Client, r *resty.Response) (time.Duration, error) {
	retryIntervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	attempt := r.Request.Attempt

	if attempt <= len(retryIntervals) {
		logger.Instance.Info(fmt.Sprintf("Попытка %d, ждем %v перед следующим запросом...\n", attempt, retryIntervals[attempt-1]))
		return retryIntervals[attempt-1], nil
	}

	return 0, fmt.Errorf("превышено количество попыток")
}

var restyClient = resty.New().
	SetRetryCount(3).
	SetTimeout(10 * time.Second).
	SetRetryWaitTime(1 * time.Second).
	SetRetryMaxWaitTime(5 * time.Second).
	SetRetryAfter(retryAfterFunc).
	AddRetryCondition(func(r *resty.Response, err error) bool {
		return err != nil // Повторяем только при сетевых ошибках
	})

func sendPostRequest(url string, body []byte, headers map[string]string) (*resty.Response, error) {
	//ЗАПОМНИТЬ: resty автоматически добавляет заголовок "Accept-Encoding", "gzip" и распаковывает ответ если он пришел в gzip
	return restyClient.R().
		SetBody(body).
		SetHeaders(headers).
		Post(url)
}

func sendWithCompression(url string, data interface{}, headers map[string]string) (*resty.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		logger.Instance.Warnw("sendWithCompression", "json.Marshal err", err)
		return &resty.Response{}, err
	}

	headers["Content-Encoding"] = "gzip"

	compressedBody, err := gzip.Compress(body)
	if err != nil {
		logger.Instance.Warnw("sendWithCompressionAndRequestId", "compress err", err)
		return &resty.Response{}, err
	}

	return sendPostRequest(url, compressedBody, headers)
}

func SendMetricBatch(domain string, data []consts.MetricData) error {
	requestID := uuid.New().String()

	headers := map[string]string{
		"Content-Type": "application/json",
		"x-request-id": requestID,
	}

	url := fmt.Sprintf("http://%s/updates/", domain)

	resp, err := sendWithCompression(url, data, headers)

	if resp.StatusCode() == 200 {
		logger.Instance.Info("Метрики успешно отправлены")
	}

	var loggerFunc = logger.Instance.Infow
	if err != nil {
		loggerFunc = logger.Instance.Warnw
	}

	loggerFunc("SendMetric Response", "requestID", requestID, "url", url, "status", resp.Status(), "body", resp.Body(), "err", err)

	return nil
}

func SendMetric(domain string, metricType consts.Metric, name string, value interface{}) error {
	requestData, err := converter.GenerateMetricData(metricType, name, value)
	if err != nil {
		logger.Instance.Warnw("SendMetric", "GenerateMetricData", err)
		return err
	}

	requestID := uuid.New().String()

	headers := map[string]string{
		"Content-Type": "application/json",
		"x-request-id": requestID,
	}

	url := fmt.Sprintf("http://%s/update/", domain)

	resp, err := sendWithCompression(url, requestData, headers)

	if resp.StatusCode() == 200 {
		logger.Instance.Info("Метрики успешно отправлены")
	}

	if err != nil {
		logger.Instance.Errorw("SendMetric", "не удалось выполнить запрос", err)
	}

	respBody := resp.Body()

	logger.Instance.Infow("SendMetric Response", "requestID", requestID, "url", url, "status", resp.Status(), "body", respBody)

	return nil
}

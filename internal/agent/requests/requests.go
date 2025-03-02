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

func SendMetricBatch(domain string, data []consts.MetricData) error {
	body, err := json.Marshal(data)
	if err != nil {
		logger.Instance.Warnw("SendMetric", "json.Marshal err", err)
		return err
	}

	compressedBody, err := gzip.Compress(body)
	if err != nil {
		logger.Instance.Warnw("SendMetric", "compress err", err)
		return err
	}

	url := fmt.Sprintf("http://%s/updates/", domain)

	requestID := uuid.New().String()

	//ЗАПОМНИТЬ: resty автоматически добавляет заголовок "Accept-Encoding", "gzip" и распаковывает ответ если он пришел в gzip
	resp, err := restyClient.R().
		SetBody(compressedBody).
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json").
		SetHeader("x-request-id", requestID).
		Post(url)

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

	//ЗАПОМНИТЬ: resty автоматически добавляет заголовок "Accept-Encoding", "gzip" и распаковывает ответ если он пришел в gzip
	resp, err := restyClient.R().
		SetBody(compressedBody).
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json").
		Post(url)

	if resp.StatusCode() == 200 {
		logger.Instance.Info("Метрики успешно отправлены")
	}

	if err != nil {
		logger.Instance.Errorw("SendMetric", "не удалось выполнить запрос", err)
	}

	respBody := resp.Body()

	logger.Instance.Infow("SendMetric Response", "url", url, "status", resp.Status(), "body", respBody)

	return nil
}

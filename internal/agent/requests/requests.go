package requests

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/converter"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/gzip"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

var retryAfterFunc = func(retryIntervals []time.Duration) func(c *resty.Client, r *resty.Response) (time.Duration, error) {

	return func(c *resty.Client, r *resty.Response) (time.Duration, error) {
		attempt := r.Request.Attempt

		if attempt <= len(retryIntervals) {
			logger.Instance.Info(fmt.Sprintf("Попытка %d, ждем %v перед следующим запросом...\n", attempt, retryIntervals[attempt-1]))
			return retryIntervals[attempt-1], nil
		}

		return 0, fmt.Errorf("превышено количество попыток")
	}
}

type Requester struct {
	client  *resty.Client
	hashKey string
	host    string
}

func NewRequester(cfg setup.AgentStartupValues) *Requester {
	var restyClient = resty.New().
		SetRetryCount(len(cfg.RetryIntervals)).
		SetTimeout(cfg.RequestWaitTimeout).
		SetRetryWaitTime(minRetryDuration(cfg.RetryIntervals)).
		SetRetryMaxWaitTime(maxRetryDuration(cfg.RetryIntervals)).
		SetRetryAfter(retryAfterFunc(cfg.RetryIntervals)).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return err != nil // Повторяем только при сетевых ошибках
		})

	return &Requester{
		client:  restyClient,
		hashKey: cfg.HashKey,
		host:    cfg.Host,
	}
}

func minRetryDuration(retryIntervals []time.Duration) time.Duration {
	minimun := retryIntervals[0]
	for _, interval := range retryIntervals {
		if interval < minimun {
			minimun = interval
		}
	}
	return minimun
}

func maxRetryDuration(retryIntervals []time.Duration) time.Duration {
	maximum := retryIntervals[0]
	for _, interval := range retryIntervals {
		if interval > maximum {
			maximum = interval
		}
	}
	return maximum
}

func (r *Requester) sendPostRequest(url string, body []byte, headers map[string]string) (*resty.Response, error) {
	//ЗАПОМНИТЬ: resty автоматически добавляет заголовок "Accept-Encoding", "gzip" и распаковывает ответ если он пришел в gzip
	return r.client.R().
		SetBody(body).
		SetHeaders(headers).
		Post(url)
}

func (r *Requester) sendWithCompression(url string, data interface{}, headers map[string]string, hashKey string) (*resty.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		logger.Instance.Warnw("sendWithCompression", "json.Marshal err", err)
		return &resty.Response{}, err
	}

	if hashKey != "" {
		h := hmac.New(sha256.New, []byte(hashKey))
		h.Write(body)
		headers["HashSHA256"] = hex.EncodeToString(h.Sum(nil))
	}

	headers["Content-Encoding"] = "gzip"

	compressedBody, err := gzip.Compress(body)
	if err != nil {
		logger.Instance.Warnw("sendWithCompressionAndRequestId", "compress err", err)
		return &resty.Response{}, err
	}

	return r.sendPostRequest(url, compressedBody, headers)
}

func (r *Requester) SendMetricBatch(data []consts.MetricData) error {
	requestID := uuid.New().String()

	headers := map[string]string{
		"Content-Type": "application/json",
		"x-request-id": requestID,
	}

	url := fmt.Sprintf("http://%s/updates/", r.host)

	resp, err := r.sendWithCompression(url, data, headers, r.hashKey)

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

func (r *Requester) SendMetric(metricType consts.Metric, name string, value interface{}) error {
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

	url := fmt.Sprintf("http://%s/update/", r.host)

	resp, err := r.sendWithCompression(url, requestData, headers, r.hashKey)

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

package requests

import (
	"encoding/json"
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/processData"
	"net"
	"time"

	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/converter"
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
	client              *resty.Client
	hashKey             string
	host                string
	cryptoPublicKeyPath string
	outboundIP          net.IP
}

func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
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

	ip, err := getOutboundIP()
	if err != nil {
		logger.Instance.Warn("error when try to get outbound IP address", err)
	}
	logger.Instance.Infof("Outbound IP: %s", ip)

	return &Requester{
		client:              restyClient,
		hashKey:             cfg.HashKey,
		host:                cfg.Host,
		cryptoPublicKeyPath: cfg.CryptoPublicKeyPath,
		outboundIP:          ip,
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

func (r *Requester) sendWithProcessedData(url string, data interface{}, headers map[string]string) (*resty.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		logger.Instance.Warnw("sendWithProcessedData", "json.Marshal err", err)
		return &resty.Response{}, err
	}

	processors := []processData.DataProcessor{
		&processData.Sha256Processor{HashKey: r.hashKey},
		&processData.XRealIpProcessor{OutboundIP: r.outboundIP},
		&processData.GZipProcessor{},
		&processData.EncryptBodyProcessor{CryptoPublicKeyPath: r.cryptoPublicKeyPath},
	}

	chainProcessor := &processData.ChainProcessor{
		Processors: processors,
	}

	processedBody, headers, err := chainProcessor.Process(body, headers)

	logger.Instance.Infow("sendWithProcessedData", "processedBody", processedBody, "headers", headers)

	if err != nil {
		logger.Instance.Warnw("chainProcessor.Process", "process err", err)
		return &resty.Response{}, err
	}

	return r.sendPostRequest(url, processedBody, headers)
}

func (r *Requester) SendMetricBatch(data []consts.MetricData) error {
	requestID := uuid.New().String()

	headers := map[string]string{
		"Content-Type": "application/json",
		"x-request-id": requestID,
	}

	url := fmt.Sprintf("http://%s/updates/", r.host)

	resp, err := r.sendWithProcessedData(url, data, headers)

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

	if r.outboundIP != nil {
		headers["X-real-ip"] = r.outboundIP.String()
	}

	url := fmt.Sprintf("http://%s/update/", r.host)

	resp, err := r.sendWithProcessedData(url, requestData, headers)

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

package requests

import (
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/consts"
	"net/http"
)

func SendMetric(metricType consts.Metric, name string, value string) error {
	domain := "localhost:8080"
	url := fmt.Sprintf("http://%s/update/%s/%s/%s", domain, metricType, name, value)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("не удалось создать реквест: %s", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос: %s", err)
	}

	defer resp.Body.Close()

	fmt.Println("Метрики успешно отправлены", err)

	return nil
}

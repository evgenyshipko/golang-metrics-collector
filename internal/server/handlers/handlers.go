package handlers

import (
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/parser"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
	"net/http"
)

func PostMetric(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Error(res, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	if parser.IsNameMissed(req.URL.Path) {
		http.Error(res, "Отсутствует имя метрики", http.StatusNotFound)
		return
	}

	data, err := parser.ParseURLPath(req.URL.Path)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = storage.STORAGE.Set(data.MetricType, data.Name, data.Value)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(storage.STORAGE)

	res.Write([]byte("Метрика записана успешно!"))
}

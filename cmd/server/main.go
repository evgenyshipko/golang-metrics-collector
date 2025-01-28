package main

import (
	"fmt"
	"github.com/evgenyshipko/golang-metrics-collector/internal/parser"
	"github.com/evgenyshipko/golang-metrics-collector/internal/storage"
	"net/http"
)

var STORAGE = storage.NewMemStorage()

func postMetric(res http.ResponseWriter, req *http.Request) {

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

	err = STORAGE.Set(data.MetricType, data.Name, data.Value)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(STORAGE)

	res.Write([]byte("Метрика записана успешно!"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, postMetric)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

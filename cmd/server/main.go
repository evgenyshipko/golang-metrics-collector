package main

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/server"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/setup"
	"net/http"
)

func main() {
	defer logger.Sync()

	s := server.Setup()

	values, err := setup.GetStartupValues()
	if err != nil {
		logger.Instance.Fatalw("Аргументы не прошли валидацию", err)
	}

	err = http.ListenAndServe(values.Host, s.Routes())
	if err != nil {
		panic(err)
	}
}

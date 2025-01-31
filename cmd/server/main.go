package main

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/router"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/setup"
	"net/http"
)

func main() {
	r := router.MakeChiRouter()

	values := setup.GetStartupValues()

	err := http.ListenAndServe(values.Host, r)
	if err != nil {
		panic(err)
	}
}

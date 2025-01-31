package main

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/router"
	"net/http"
)

func main() {
	r := router.MakeChiRouter()

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}

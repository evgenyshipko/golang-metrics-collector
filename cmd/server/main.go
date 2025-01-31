package main

import (
	"flag"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/router"
	"net/http"
)

func main() {
	r := router.MakeChiRouter()

	host := flag.String("a", "localhost:8080", "input host with port")
	flag.Parse()

	err := http.ListenAndServe(*host, r)
	if err != nil {
		panic(err)
	}
}

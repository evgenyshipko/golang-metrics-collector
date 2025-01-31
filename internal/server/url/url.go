package url

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/consts"
	"github.com/go-chi/chi"
	"net/http"
)

func URLParam(r *http.Request, key consts.UrlParam) string {
	return chi.URLParam(r, string(key))
}

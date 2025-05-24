package url

import (
	"net/http"

	"github.com/evgenyshipko/golang-metrics-collector/internal/common/consts"
	"github.com/go-chi/chi"
)

func MyURLParam(r *http.Request, key consts.URLParam) string {
	return chi.URLParam(r, string(key))
}

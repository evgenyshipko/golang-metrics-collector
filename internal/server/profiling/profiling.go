package profiling

import (
	"github.com/go-chi/chi"
	"net/http"
	"net/http/pprof"
)

func PprofHandlers() http.Handler {
	r := chi.NewRouter()

	// Регистрируем все стандартные обработчики pprof
	r.HandleFunc("/", pprof.Index)
	r.HandleFunc("/cmdline", pprof.Cmdline)
	r.HandleFunc("/profile", pprof.Profile)
	r.HandleFunc("/symbol", pprof.Symbol)
	r.HandleFunc("/trace", pprof.Trace)

	// Для heap и других профилей
	r.Handle("/goroutine", pprof.Handler("goroutine"))
	r.Handle("/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/heap", pprof.Handler("heap"))
	r.Handle("/block", pprof.Handler("block"))
	r.Handle("/mutex", pprof.Handler("mutex"))
	r.Handle("/allocs", pprof.Handler("allocs"))

	return r
}

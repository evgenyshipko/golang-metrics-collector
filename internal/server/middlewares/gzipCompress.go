package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipCompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			logger.Instance.Warnw("GzipCompress", "Error creating gzip writer", err)
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

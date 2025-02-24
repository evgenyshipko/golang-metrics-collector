package logging

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"net/http"
	"time"
)

func LoggingHandlers(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &ResponseData{
			Status: 0,
			Size:   0,
		}
		lw := LoggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			ResponseData:   responseData,
		}
		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)

		loggerFunc := logger.Instance.Infow
		if responseData.Status >= 400 {
			loggerFunc = logger.Instance.Warnw
		}

		loggerFunc("Request",
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.Status,
			"duration", duration,
			"size", responseData.Size,
			"responseBody", lw.Body.String(),
		)
	}
	return http.HandlerFunc(logFn)
}

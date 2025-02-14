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

		// Если статус 4xx или 5xx — сохраняем тело ошибки
		if responseData.Status >= 400 {
			responseData.Error = lw.Body.String()
		}

		logger.Instance.Infow("Request",
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.Status,
			"duration", duration,
			"size", responseData.Size,
			"error", responseData.Error,
		)
	}
	return http.HandlerFunc(logFn)
}

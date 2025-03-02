package logging

import (
	"bytes"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/go-chi/chi/middleware"
	"io"
	"net/http"
	"time"
)

func LoggingHandlers(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Получаем requestID из контекста (middleware chi)
		requestID := middleware.GetReqID(r.Context())

		// Читаем тело запроса
		var requestBody string
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				// Восстанавливаем r.Body, чтобы обработчики могли его использовать
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		loggerFunc := logger.Instance.Infow

		loggerFunc("Request",
			"requestId", requestID,
			"uri", r.RequestURI,
			"method", r.Method,
			"body", requestBody,
		)

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

		if responseData.Status >= 400 {
			loggerFunc = logger.Instance.Warnw
		}

		loggerFunc("Response",
			"requestId", requestID,
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.Status,
			"duration", duration,
			"size", responseData.Size,
			"body", lw.Body.String(),
		)
	}
	return http.HandlerFunc(logFn)
}

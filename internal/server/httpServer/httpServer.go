package httpServer

import (
	"context"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"net/http"
	"time"
)

type HttpServer struct {
	server http.Server
}

func NewHttpServer(host string, handler http.Handler) *HttpServer {
	return &HttpServer{
		server: http.Server{Addr: host, Handler: handler},
	}
}

func (s *HttpServer) Start() error {
	logger.Instance.Infow("SERVER STARTED!")
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Instance.Fatalw("httpServer.ListenAndServe", "Ошибка запуска сервера", err)
		return err
	}
	return nil
}

func (s *HttpServer) Stop() error {
	// Создаём контекст с таймаутом для корректного завершения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		logger.Instance.Warnw("httpServer.Shutdown", "Ошибка завершения сервера:", err)
		return err
	}
	return nil
}

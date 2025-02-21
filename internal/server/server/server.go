package server

import (
	"context"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/files"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares/logging"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"time"
)

type CustomServer struct {
	http.Server
	router *chi.Mux
	store  *storage.MemStorage
	config *setup.ServerStartupValues
}

func NewServer(router *chi.Mux, store *storage.MemStorage, config *setup.ServerStartupValues) *CustomServer {
	s := &CustomServer{
		Server: http.Server{Addr: config.Host, Handler: router},
		router: router,
		store:  store,
		config: config,
	}
	s.routes()
	return s
}

func (s *CustomServer) Routes() *chi.Mux {
	return s.router
}

func (s *CustomServer) GetStoreData() *storage.MemStorageData {
	return s.store.GetAll()
}

func Create(config *setup.ServerStartupValues) *CustomServer {
	router := chi.NewRouter()
	// TODO: логгировать RequestId
	router.Use(middleware.RequestID)

	router.Use(middlewares.GzipDecompress)

	router.Use(middlewares.GzipCompress)

	router.Use(logging.LoggingHandlers)

	store := storage.NewMemStorage()

	if config.Restore {
		files.ReadFromFile(config.FileStoragePath, store)
	}

	server := NewServer(router, store, config)
	return server
}

func (s *CustomServer) Start() {
	logger.Instance.Infow("SERVER STARTED!")
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		logger.Instance.Fatalw("httpServer.ListenAndServe", "Ошибка запуска сервера", err)
	}
}

func (s *CustomServer) ShutDown() {
	logger.Instance.Info("Завершение сервера...")

	// Создаём контекст с таймаутом для корректного завершения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Server.Shutdown(ctx); err != nil {
		logger.Instance.Warnw("httpServer.Shutdown", "Ошибка завершения сервера:", err)
	}

	files.WriteToFile(s.config.FileStoragePath, s.store.GetAll())

	logger.Instance.Info("Сервер успешно завершён")
}

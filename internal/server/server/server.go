package server

import (
	"context"
	"database/sql"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	db2 "github.com/evgenyshipko/golang-metrics-collector/internal/server/db"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/files"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares/logging"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/services"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
	"github.com/go-chi/chi"
	"net/http"
	"time"
)

type CustomServer struct {
	server  http.Server
	router  *chi.Mux
	store   storage.Storage
	config  *setup.ServerStartupValues
	service services.Service
	db      *sql.DB
}

func NewCustomServer(router *chi.Mux, store storage.Storage, config *setup.ServerStartupValues, service services.Service, db *sql.DB) *CustomServer {
	s := &CustomServer{
		server:  http.Server{Addr: config.Host, Handler: router},
		router:  router,
		store:   store,
		config:  config,
		service: service,
		db:      db,
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

	router.Use(middlewares.GzipDecompress)

	router.Use(middlewares.GzipCompress)

	router.Use(logging.LoggingHandlers)

	store := storage.NewMemStorage()

	if config.Restore {
		files.ReadFromFile(config.FileStoragePath, store)
	}

	service := services.NewMetricService(store, config.StoreInterval, config.FileStoragePath)

	db := db2.ConnectToDb(config.DatabaseDSN)

	server := NewCustomServer(router, store, config, service, db)
	return server
}

func (s *CustomServer) Start() {
	logger.Instance.Infow("SERVER STARTED!")
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Instance.Fatalw("httpServer.ListenAndServe", "Ошибка запуска сервера", err)
	}
}

func (s *CustomServer) ShutDown() {
	logger.Instance.Info("Завершение сервера...")

	// Создаём контекст с таймаутом для корректного завершения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		logger.Instance.Warnw("httpServer.Shutdown", "Ошибка завершения сервера:", err)
	}

	files.WriteToFile(s.config.FileStoragePath, s.store.GetAll())

	logger.Instance.Info("Сервер успешно завершён")
}

func (s *CustomServer) GetDB() *sql.DB {
	return s.db
}

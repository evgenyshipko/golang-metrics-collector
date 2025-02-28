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
	"github.com/go-chi/chi/middleware"
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

func (s *CustomServer) GetStoreData() (*storage.StorageData, error) {
	return s.store.GetAll()
}

func Create(config *setup.ServerStartupValues) *CustomServer {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)

	router.Use(middlewares.GzipDecompress)

	router.Use(middlewares.GzipCompress)

	router.Use(logging.LoggingHandlers)

	db := db2.ConnectToDB(config.DatabaseDSN)

	var store storage.Storage
	if config.DatabaseDSN != "" {
		store = storage.NewSQLStorage(db)
	} else {
		store = storage.NewMemStorage()
	}

	if config.Restore {
		files.ReadFromFile(config.FileStoragePath, store)
	}

	service := services.NewMetricService(store, config.StoreInterval, config.FileStoragePath)

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

	data, err := s.store.GetAll()
	if err != nil {
		logger.Instance.Warnw("ShutDown", "ошибка получения данных", err)
		return
	}

	files.WriteToFile(s.config.FileStoragePath, data)

	logger.Instance.Info("Сервер успешно завершён")
}

func (s *CustomServer) GetDB() *sql.DB {
	return s.db
}

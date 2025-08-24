package server

import (
	"context"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/files"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/grpcServer"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/profiling"
	_ "net/http/pprof"

	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/httpserver"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares/logging"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/services"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server interface {
	Start() error
	Stop() error
}

func CreateServer(router *chi.Mux, config *setup.ServerStartupValues, service services.Service) Server {
	switch config.Protocol {
	case "http":
		return httpserver.NewHTTPServer(config.Host, router)
	default:
		return grpcServer.CreateGrpcServer(service, *config)
	}
}

type CustomServer struct {
	server  Server
	router  *chi.Mux
	store   storage.Storage
	config  *setup.ServerStartupValues
	service services.Service
}

func NewCustomServer(router *chi.Mux, store storage.Storage, config *setup.ServerStartupValues) *CustomServer {
	service := services.NewMetricService(store, config.StoreInterval, config.FileStoragePath)
	s := &CustomServer{
		server:  CreateServer(router, config, service),
		router:  router,
		store:   store,
		config:  config,
		service: service,
	}
	s.routes()
	return s
}

func (s *CustomServer) GetStoreData() (*storage.StorageData, error) {
	return s.store.GetAll(context.Background())
}

func Create(config *setup.ServerStartupValues, store storage.Storage) *CustomServer {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)

	router.Use(middlewares.TrustedIpMiddleware(config.TrustedSubnet))

	router.Use(middlewares.DecryptMiddleware(config.CryptoPrivateKeyPath))

	router.Use(middlewares.GzipDecompress)

	router.Use(middlewares.GzipCompress)

	router.Use(logging.LoggingHandlers)

	router.Use(middlewares.ValidateSHA256(config.HashKey))

	router.Mount("/debug/pprof", profiling.PprofHandlers())

	server := NewCustomServer(router, store, config)
	return server
}

func (s *CustomServer) Start() {
	err := s.server.Start()
	if err != nil {
		logger.Instance.Warn("Failed to start server")
	}
}

func (s *CustomServer) ShutDown() {
	logger.Instance.Info("Завершение сервера...")

	err := s.server.Stop()
	if err != nil {
		logger.Instance.Warnw("CustomServer.Shutdown", "Ошибка завершения сервера Stop()", err)
	}

	data, err := s.store.GetAll(context.Background())
	if err != nil {
		logger.Instance.Warnw("ShutDown", "ошибка получения данных", err)
		return
	}

	err = files.WriteToFileWithRetry(s.config.FileStoragePath, data, s.config.RetryIntervals)
	if err != nil {
		logger.Instance.Warnw("Не удалось записать в файл по завершению сервера", "Ошибка", err)
	}

	logger.Instance.Info("Сервер успешно завершён")
}

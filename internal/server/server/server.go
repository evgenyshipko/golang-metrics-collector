package server

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares/logging"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	router *chi.Mux
	store  *storage.MemStorage
}

func NewServer(router *chi.Mux, store *storage.MemStorage) *Server {
	s := &Server{
		router: router,
		store:  store,
	}
	s.routes()
	return s
}

func (s *Server) Routes() *chi.Mux {
	return s.router
}

func Setup() *Server {
	router := chi.NewRouter()

	// TODO: логгировать RequestId
	router.Use(middleware.RequestID)

	router.Use(middlewares.GzipDecompress)

	router.Use(middlewares.GzipCompress)

	router.Use(logging.LoggingHandlers)

	store := storage.NewMemStorage()
	server := NewServer(router, store)
	return server
}

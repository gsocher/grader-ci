package server

import (
	"net/http"

	"github.com/dpolansky/ci/server/route"
	"github.com/dpolansky/ci/server/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Router            *mux.Router
	Builder           service.Builder
	RepositoryService service.RepositoryService
}

// New initializes a server with its dependencies and registers its routes.
func New(builder service.Builder, repositoryService service.RepositoryService) (*Server, error) {
	s := &Server{
		Router:            mux.NewRouter(),
		Builder:           builder,
		RepositoryService: repositoryService,
	}

	// register routes
	s.registerRoutes()

	return s, nil
}

func (s *Server) Serve() {
	serv := &http.Server{
		Addr:    "localhost:8081",
		Handler: s.Router,
	}

	logrus.Infof("Starting server on %v", serv.Addr)
	go s.Builder.ListenForUpdates()

	logrus.Fatalf("Server shut down: %v\n", serv.ListenAndServe())
}

func (s *Server) registerRoutes() {
	// api routes
	route.RegisterGithubWebhookRoutes(s.Router, s.Builder)
	route.RegisterBuildStatusRoutes(s.Router, s.Builder)
	route.RegisterRepositoryRoutes(s.Router, s.RepositoryService)

	// frontend routes
	route.RegisterBuildStatusFrontendRoutes(s.Router, s.Builder)
	route.RegisterRepositoryFrontendRoutes(s.Router, s.RepositoryService)
}

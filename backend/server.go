package backend

import (
	"net/http"

	"github.com/dpolansky/ci/backend/route"
	"github.com/dpolansky/ci/backend/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Router               *mux.Router
	BuildReadWriter      service.BuildReadWriter
	BuildRunner          service.BuildRunner
	RepositoryReadWriter service.RepositoryReadWriter
}

// New initializes a server with its dependencies and registers its routes.
func New(build service.BuildReadWriter, run service.BuildRunner, rep service.RepositoryReadWriter) (*Server, error) {
	s := &Server{
		Router:               mux.NewRouter(),
		BuildReadWriter:      build,
		RepositoryReadWriter: rep,
		BuildRunner:          run,
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
	go s.BuildRunner.ListenForUpdates()

	logrus.Fatalf("Server shut down: %v\n", serv.ListenAndServe())
}

func (s *Server) registerRoutes() {
	// api routes
	route.RegisterGithubWebhookRoutes(s.Router, s.BuildRunner, s.RepositoryReadWriter)
	route.RegisterBuildRoutes(s.Router, s.BuildReadWriter)
	route.RegisterRepositoryRoutes(s.Router, s.RepositoryReadWriter)
}

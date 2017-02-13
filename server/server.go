package server

import (
	"net/http"

	"github.com/dpolansky/ci/server/amqp"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Router       *mux.Router
	BuildService BuildService
}

// New initializes a server with its dependencies and registers its routes.
func New(amqpClient amqp.ReadWriter, buildService BuildService) (*Server, error) {
	s := &Server{
		Router:       mux.NewRouter(),
		BuildService: buildService,
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
	go s.BuildService.ListenForUpdates()

	logrus.Fatalf("Server shut down: %v\n", serv.ListenAndServe())
}

func (s *Server) registerRoutes() {
	s.registerGithubWebhookRoutes()
	s.registerBuildStatusRoutes()
}

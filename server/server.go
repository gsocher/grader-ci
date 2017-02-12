package server

import (
	"fmt"
	"net/http"

	"github.com/dpolansky/ci/server/amqp"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Router       *mux.Router
	BuildService BuildService
}

func New() (*Server, error) {
	// init amqp client
	client, err := amqp.NewClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize amqp client: %v", err)
	}

	// init build service
	buildService := NewBuildService(client)

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
	logrus.Fatalf("Server shut down: %v\n", serv.ListenAndServe())
}

func (s *Server) registerRoutes() {
	s.registerGithubWebhookRoutes()
}

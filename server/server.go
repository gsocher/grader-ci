package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Router *mux.Router
}

func New() *Server {
	s := &Server{
		Router: mux.NewRouter(),
	}

	// register routes
	s.registerRoutes()

	return s
}

func (s *Server) Serve() {
	serv := &http.Server{
		Addr:    "localhost:8080",
		Handler: s.Router,
	}

	logrus.Infof("Starting server on %v", serv.Addr)
	logrus.Fatalf("Server shut down: %v\n", serv.ListenAndServe())
}

func (s *Server) registerRoutes() {
	s.registerGithubWebhookRoutes()
}

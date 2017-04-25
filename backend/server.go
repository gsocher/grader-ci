package backend

import (
	"net/http"

	"github.com/dpolansky/grader-ci/backend/route"
	"github.com/dpolansky/grader-ci/backend/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Router               *mux.Router
	BuildReadWriter      service.BuildReadWriter
	BuildRunner          service.BuildRunner
	RepositoryReadWriter service.RepositoryReadWriter
	TestBindReadWriter   service.TestBindReadWriter
}

// New initializes a server with its dependencies and registers its routes.
func New(build service.BuildReadWriter, run service.BuildRunner, rep service.RepositoryReadWriter, bind service.TestBindReadWriter) (*Server, error) {
	s := &Server{
		Router:               mux.NewRouter(),
		BuildReadWriter:      build,
		RepositoryReadWriter: rep,
		BuildRunner:          run,
		TestBindReadWriter:   bind,
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
	route.RegisterGithubWebhookRoutes(s.Router, s.BuildRunner, s.RepositoryReadWriter, s.TestBindReadWriter)
	route.RegisterBuildAPIRoutes(s.Router, s.BuildReadWriter)
	route.RegisterRepositoryAPIRoutes(s.Router, s.RepositoryReadWriter)
	route.RegisterBindAPIRoutes(s.Router, s.TestBindReadWriter)

	// frontend routes
	route.RegisterRepositoryFrontendRoutes(s.Router, s.RepositoryReadWriter)
	route.RegisterBuildFrontendRoutes(s.Router, s.BuildReadWriter, s.RepositoryReadWriter)
	route.RegisterBindFrontendRoutes(s.Router, s.TestBindReadWriter)

	// assets route
	route.RegisterAssetsRoute(s.Router)
}

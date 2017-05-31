package backend

import (
	"net/http"

	"github.com/dpolansky/grader-ci/pkg/backend/route"
	"github.com/dpolansky/grader-ci/pkg/backend/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server interface {
	Run()
}

type server struct {
	config *ServerConfig
}

type ServerConfig struct {
	Router               *mux.Router
	BuildService         service.BuildService
	BuildMessageService  service.BuildMessageService
	RepositoryService    service.RepositoryService
	TestBindService      service.TestBindService
	GithubWebhookService service.GithubWebhookService
}

// New initializes a server with its dependencies and registers its routes.
func New(config *ServerConfig) (Server, error) {
	if config.Router == nil {
		config.Router = mux.NewRouter()
	}

	s := &server{config: config}

	// register routes
	s.registerRoutes()

	return s, nil
}

func (s *server) Run() {
	serv := &http.Server{
		Addr:    "localhost:8081",
		Handler: s.config.Router,
	}

	logrus.Infof("Starting server on %v", serv.Addr)
	go s.config.BuildMessageService.ListenForBuildMessages(nil)

	logrus.Fatalf("Server shut down: %v\n", serv.ListenAndServe())
}

func (s *server) registerRoutes() {
	// api routes
	route.RegisterGithubWebhookRoutes(s.config.Router, s.config.GithubWebhookService)
	route.RegisterBuildAPIRoutes(s.config.Router, s.config.BuildService)
	route.RegisterRepositoryAPIRoutes(s.config.Router, s.config.RepositoryService)
	route.RegisterBindAPIRoutes(s.config.Router, s.config.TestBindService)

	// frontend routes
	route.RegisterRepositoryFrontendRoutes(s.config.Router, s.config.RepositoryService)
	route.RegisterBuildFrontendRoutes(s.config.Router, s.config.BuildService, s.config.RepositoryService)
	route.RegisterBindFrontendRoutes(s.config.Router, s.config.TestBindService)

	// assets route
	route.RegisterAssetsRoute(s.config.Router)
}

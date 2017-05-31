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
	BuildService         service.BuildService
	BuildMessageService  service.BuildMessageService
	RepositoryService    service.RepositoryService
	TestBindService      service.TestBindService
	GithubWebhookService service.GithubWebhookService
}

// New initializes a server with its dependencies and registers its routes.
func New(build service.BuildService, msg service.BuildMessageService, rep service.RepositoryService, bind service.TestBindService, github service.GithubWebhookService) (*Server, error) {
	s := &Server{
		Router:               mux.NewRouter(),
		BuildService:         build,
		RepositoryService:    rep,
		BuildMessageService:  msg,
		TestBindService:      bind,
		GithubWebhookService: github,
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
	go s.BuildMessageService.ListenForBuildMessages()

	logrus.Fatalf("Server shut down: %v\n", serv.ListenAndServe())
}

func (s *Server) registerRoutes() {
	// api routes
	route.RegisterGithubWebhookRoutes(s.Router, s.GithubWebhookService)
	route.RegisterBuildAPIRoutes(s.Router, s.BuildService)
	route.RegisterRepositoryAPIRoutes(s.Router, s.RepositoryService)
	route.RegisterBindAPIRoutes(s.Router, s.TestBindService)

	// frontend routes
	route.RegisterRepositoryFrontendRoutes(s.Router, s.RepositoryService)
	route.RegisterBuildFrontendRoutes(s.Router, s.BuildService, s.RepositoryService)
	route.RegisterBindFrontendRoutes(s.Router, s.TestBindService)

	// assets route
	route.RegisterAssetsRoute(s.Router)
}

package server

import (
	"html/template"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/dpolansky/ci/model"
)

func TestStatusTemplate(t *testing.T) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/status.html"))
		statuses := []model.BuildStatus{
			{
				ID:         1,
				LastUpdate: time.Now(),
				CloneURL:   "github.com/docker/docker",
				Branch:     "master",
				Status:     model.StatusBuildPassed,
				Log:        "log output...",
			},
		}

		tmpl.Execute(w, struct{ Statuses []model.BuildStatus }{statuses})
	})

	log.Printf("listening...")
	http.ListenAndServe(":8090", nil)
}

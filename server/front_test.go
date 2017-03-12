package server

import (
	"fmt"
	"html/template"
	"net/http"
	"testing"
	"time"

	"github.com/dpolansky/ci/model"
)

func TestStatusTemplate(t *testing.T) {

	http.HandleFunc("/detail", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("got detail req\n")
		tmpl := template.Must(template.ParseFiles("templates/status_detail.html"))
		status := model.BuildStatus{
			ID:         1,
			LastUpdate: time.Now(),
			CloneURL:   "github.com/docker/docker",
			Branch:     "master",
			Status:     model.StatusBuildPassed,
			Log:        "log output...",
		}

		tmpl.Execute(w, status)
	})

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	tmpl := template.Must(template.ParseFiles("templates/status_list.html"))
	// 	statuses := []model.BuildStatus{
	// 		{
	// 			ID:         1,
	// 			LastUpdate: time.Now(),
	// 			CloneURL:   "github.com/docker/docker",
	// 			Branch:     "master",
	// 			Status:     model.StatusBuildPassed,
	// 			Log:        "log output...",
	// 		},
	// 	}

	// 	tmpl.Execute(w, struct{ Statuses []model.BuildStatus }{statuses})
	// })

	fmt.Printf("listening...\n")
	http.ListenAndServe(":8090", nil)
}

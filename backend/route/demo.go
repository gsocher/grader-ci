package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"os"

	"github.com/dpolansky/grader-ci/backend/service"
	"github.com/dpolansky/grader-ci/model"
)

type Repo struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Owner    struct {
		Login     string `json:"login"`
		ID        int    `json:"id"`
		AvatarURL string `json:"avatar_url"`
	} `json:"owner"`
}

func runDemoHTTPHandler(run service.BuildRunner, rep service.RepositoryReadWriter, bind service.TestBindReader) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		b, err := ioutil.ReadFile(filepath.Join(os.Getenv("GOPATH"), "src/github.com/dpolansky/grader-ci/repos.json"))
		if err != nil {
			writeError(rw, http.StatusInternalServerError, err)
			return
		}

		if err = createDemoGraderRepository(rep); err != nil {
			writeError(rw, http.StatusInternalServerError, err)
			return
		}

		var repos []Repo
		json.Unmarshal(b, &repos)

		demoAPIKey := os.Getenv("DEMO_API_KEY")

		// convert to appropriate models
		for _, r := range repos {
			// create each repository
			if err = rep.UpdateRepository(&model.Repository{
				AvatarURL: r.Owner.AvatarURL,
				Owner:     r.Owner.Login,
				ID:        r.ID,
			}); err != nil {
				writeError(rw, http.StatusInternalServerError, err)
				return
			}

			_, err = run.RunBuild(&model.BuildStatus{
				Source: &model.RepositoryMetadata{
					Branch:   "master",
					CloneURL: fmt.Sprintf("https://%s@github.com/%s", demoAPIKey, r.FullName),
					ID:       r.ID,
				},
				Tested: true,
				Test: &model.RepositoryMetadata{
					Branch:   "master",
					CloneURL: "https://dpolansky@gitlab.com/dpolansky/grader-ci-demo.git",
					ID:       0,
				},
			})

			if err != nil {
				writeError(rw, http.StatusInternalServerError, err)
				return
			}
		}

		writeOk(rw, []byte{})
	}
}

func createDemoGraderRepository(rep service.RepositoryReadWriter) error {
	return rep.UpdateRepository(&model.Repository{
		AvatarURL: "https://avatars2.githubusercontent.com/u/16250555?v=3&s=460",
		Name:      "grader-ci-demo",
		Owner:     "dpolansky",
		ID:        0, // this is a github ID but we're using gitlab to store the grader repository privately so ignore this
	})
}

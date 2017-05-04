package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const apiKeyEnvVarName = "DEMO_API_KEY"

type Repo struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Owner    struct {
		Login     string `json:"login"`
		ID        int    `json:"id"`
		AvatarURL string `json:"avatar_url"`
	} `json:"owner"`
}

func main() {
	var repos []Repo

	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user/repos", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %v", os.Getenv(apiKeyEnvVarName)))

	q := req.URL.Query()
	q.Set("per_page", fmt.Sprintf("%v", 100))

	page := 1
	for {
		q.Set("page", fmt.Sprintf("%v", page))
		req.URL.RawQuery = q.Encode()

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		var temp []Repo
		json.Unmarshal(b, &temp)

		fmt.Printf("page=%v\tlen=%v\n", page, len(temp))
		if temp == nil || len(temp) == 0 {
			break
		}

		for _, r := range temp {
			fmt.Println(r.FullName)
		}

		repos = append(repos, temp...)

		page++
	}
	fmt.Printf("cleaning...\n")

	var finalRepos []Repo
	for _, r := range repos {
		if strings.Contains(r.FullName, "lab-3") {
			fmt.Println(r.FullName)
			finalRepos = append(finalRepos, r)
		}
	}

	b, _ := json.Marshal(finalRepos)
	ioutil.WriteFile("repos.json", b, 0644)
}

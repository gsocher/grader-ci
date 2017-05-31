package model

type GithubWebhookRequest struct {
	Ref        string `json:"ref"`
	Repository struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Owner struct {
			Name      string `json:"login"`
			AvatarURL string `json:"avatar_url"`
		} `json:"owner"`
	} `json:"repository"`
}

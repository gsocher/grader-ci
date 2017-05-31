package model

type Repository struct {
	ID        int    `json:"id"`
	Owner     string `json:"owner"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

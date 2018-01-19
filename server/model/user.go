package model

type User struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarURL"`
}

package models

type RootDirgodConfig struct {
	User User `json:"user"`
}

type InternalDirgodConfig struct {
	User User `json:"user"`
}

type User struct {
	Username   string `json:"username,omitempty"`
	Email      string `json:"email,omitempty"`
	Signingkey string `json:"signingkey,omitempty"`
}

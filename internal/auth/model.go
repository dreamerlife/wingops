package auth

import "time"

type User struct {
	ID           string
	Username     string
	PasswordHash string
	DisplayName  string
	Status       string
}

func (u User) Active() bool {
	return u.Status == "" || u.Status == "active"
}

type TokenPair struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type Claims struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	ExpiresAt   time.Time
}

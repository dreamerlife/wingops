package auth

import "time"

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	DisplayName  string `json:"display_name"`
	Status       string `json:"status"`
	Roles        []Role `json:"roles"`
}

func (u User) Active() bool {
	return u.Status == "" || u.Status == "active"
}

type TokenPair struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type Claims struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	DisplayName string   `json:"display_name"`
	Permissions []string `json:"permissions"`
	ExpiresAt   time.Time
}

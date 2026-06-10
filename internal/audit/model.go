package audit

import "time"

type Log struct {
	ActorID    string    `json:"actor_id"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	StatusCode int       `json:"status_code"`
	Resource   string    `json:"resource"`
	CreatedAt  time.Time `json:"created_at"`
}

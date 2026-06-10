package messaging

import (
	"errors"

	"github.com/nats-io/nats.go"
)

func Connect(url string) (*nats.Conn, error) {
	if url == "" {
		return nil, errors.New("nats url is required")
	}
	return nats.Connect(url)
}

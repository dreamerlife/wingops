package cache

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

func NewClient(addr string) (*redis.Client, error) {
	if addr == "" {
		return nil, errors.New("redis addr is required")
	}
	return redis.NewClient(&redis.Options{Addr: addr}), nil
}

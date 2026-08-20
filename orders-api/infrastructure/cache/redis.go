package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func New(url string) (*Redis, error) {
	if url == "" {
		return nil, fmt.Errorf("redis url is required")
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &Redis{
		Client: client,
	}, nil
}

func (r *Redis) Close() error {
	if r == nil || r.Client == nil {
		return nil
	}

	return r.Client.Close()
}

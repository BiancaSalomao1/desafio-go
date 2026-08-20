package cache

import (
	"context"
	"fmt"
	"time"

	"orders-api/internal/security"

	"github.com/redis/go-redis/v9"
)

type RedisTokenStore struct {
	client *redis.Client
}

func NewRedisTokenStore(
	client *redis.Client,
) *RedisTokenStore {
	return &RedisTokenStore{
		client: client,
	}
}

func (s *RedisTokenStore) Save(
	ctx context.Context,
	token string,
	ttl time.Duration,
) error {
	key := s.key(token)

	return s.client.Set(
		ctx,
		key,
		"active",
		ttl,
	).Err()
}

func (s *RedisTokenStore) Exists(
	ctx context.Context,
	token string,
) (bool, error) {
	key := s.key(token)

	count, err := s.client.Exists(
		ctx,
		key,
	).Result()

	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (s *RedisTokenStore) Delete(
	ctx context.Context,
	token string,
) error {
	key := s.key(token)

	return s.client.Del(
		ctx,
		key,
	).Err()
}

func (s *RedisTokenStore) key(
	token string,
) string {
	return fmt.Sprintf(
		"auth:token:%s",
		token,
	)
}

var _ security.TokenStore = (*RedisTokenStore)(nil)

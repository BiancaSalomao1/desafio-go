package security

import (
	"context"
	"time"
)

type TokenStore interface {
	Save(
		ctx context.Context,
		token string,
		ttl time.Duration,
	) error

	Exists(
		ctx context.Context,
		token string,
	) (bool, error)

	Delete(
		ctx context.Context,
		token string,
	) error
}

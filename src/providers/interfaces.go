package providers

import (
	"context"
	"time"
)

type KeyValServ interface {
	Get(ctx context.Context, key string) (value any, exists bool, er error)
	Set(ctx context.Context, key string, val any, expiration time.Duration) error
	Exists(ctx context.Context, key string) (bool, error)
}

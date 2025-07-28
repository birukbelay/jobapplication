package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	conf "github.com/projTemplate/goauth/src/models/config"
	"github.com/projTemplate/goauth/src/providers"
)

type RedisService struct {
	RedisClient *redis.Client
}

// Exists implements providers.KeyValServ.
func NewRedis(config *conf.KeyValConfig) (providers.KeyValServ, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.KVHost, config.KVPort),
		Password: config.KVPassword, // no password set if empty
		DB:       config.KVDbName,   // use default DB
		Username: config.KVUsername,
	})
	return &RedisService{RedisClient: rdb}, nil
}

func (r *RedisService) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if exists > 0 {
		return true, nil
	}
	return false, nil
}

// Get implements providers.KeyValServ.
func (r *RedisService) Get(ctx context.Context, key string) (value any, exists bool, er error) {
	value, err := r.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return value, false, nil
	}
	if err != nil {
		return value, false, err
	}
	return value, true, nil
}

// Set implements providers.KeyValServ.
func (r *RedisService) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	err := r.RedisClient.Set(ctx, key, val, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set value: %v", err)
	}
	return nil
}

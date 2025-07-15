package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"

	conf "github.com/projTemplate/goauth/src/models/config"
)

func NewRedis(config *conf.KeyValConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.KVHost, config.KVPort),
		Password: config.KVPassword, // no password set if empty
		DB:       config.KVDbName,   // use default DB
		Username: config.KVUsername,
	})
	return rdb, nil
}

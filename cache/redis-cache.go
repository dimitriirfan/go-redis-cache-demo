package cache

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v9"
)

type RedisCache interface {
	Set(key string, value any) error
	Get(key string) (any, error)
}

type redisCache struct {
	logger   *log.Logger
	addr     string
	password string
	db       int
	expires  time.Duration
	ctx      context.Context
}

func NewRedisCache(logger *log.Logger, addr string, password string, db int, expires time.Duration, ctx context.Context) *redisCache {
	return &redisCache{logger, addr, password, db, expires, ctx}
}

func (cache *redisCache) connectClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.addr,
		Password: cache.password,
		DB:       cache.db,
	})
}

func (cache *redisCache) Set(key string, value any) error {
	client := cache.connectClient()
	err := client.Set(cache.ctx, key, value, cache.expires).Err()

	if err != nil {
		return err
	}

	return nil

}

func (cache *redisCache) Get(key string) (any, error) {
	client := cache.connectClient()
	val, err := client.Get(cache.ctx, key).Result()

	if err != nil {
		return nil, err
	}

	return val, nil
}

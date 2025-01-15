package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, saveTime int) error
}

type RedisCache struct {
	cache *redis.Client
}

func NewRedisCache(cache *redis.Client) Cache {
	return &RedisCache{cache: cache}
}

func (r *RedisCache) Get(key string) ([]byte, error) {

	res, err := r.cache.Get(context.Background(), key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, err
		}
		err := errors.Errorf("cache.Get: %s %v", key, err)
		return nil, err
	}

	return res, nil
}

func (r *RedisCache) Set(key string, value []byte, saveTimeHour int) error {
	err := r.cache.Set(context.Background(), key, value, time.Duration(saveTimeHour)*time.Hour).Err()
	if err != nil {
		err = errors.Errorf("cache.Set: %s %v", key, err)
	}
	return err
}

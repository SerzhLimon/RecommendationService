package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/SerzhLimon/RecommendationService/internal/models"
)

type Cache interface {
	GetSong(data models.ListenSongRequest) (string, error)
	SetSong(data models.ListenSongRequest, song string, saveTime int) error
}

type RedisCache struct {
	cache *redis.Client
}

func NewRedisCache(cache *redis.Client) Cache {
	return &RedisCache{cache: cache}
}

func (r *RedisCache) GetSong(data models.ListenSongRequest) (string, error) {
	songID := strconv.FormatInt(data.SongID, 10)

	result, err := r.cache.Get(context.Background(), songID).Result()
	if err != nil {
		if err == redis.Nil {
			return "", err
		}
		err := errors.Errorf("cache.GetSong: %v", err)
		return "", err
	}

	return result, nil
}

func (r *RedisCache) SetSong(data models.ListenSongRequest, song string, saveTimeHour int) error {
	songID := strconv.FormatInt(data.SongID, 10)
	err := r.cache.Set(context.Background(), songID, song, time.Duration(saveTimeHour)*time.Hour).Err()
	if err != nil {
		err = errors.Errorf("cache.SetSong: %v", err)
	}
	return err
}

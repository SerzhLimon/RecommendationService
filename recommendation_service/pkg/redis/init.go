package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/RecommendationService/config"
)

func InitRedisClient(cfg config.Config) (*redis.Client, error) {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)

	r := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	err := r.Ping(context.Background()).Err()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"host":  cfg.Redis.Host,
			"port":  cfg.Redis.Port,
			"error": err.Error(),
		}).Error("Failed to ping Redis server")
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"host": cfg.Redis.Host,
		"port": cfg.Redis.Port,
	}).Info("Successful connection to Redis")

	return r, nil
}

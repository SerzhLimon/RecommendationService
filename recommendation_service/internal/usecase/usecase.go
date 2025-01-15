package usecase

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/RecommendationService/config"
	"github.com/SerzhLimon/RecommendationService/internal/cache"
	"github.com/SerzhLimon/RecommendationService/internal/models"
	"github.com/SerzhLimon/RecommendationService/internal/repository"
	rediscashe "github.com/SerzhLimon/RecommendationService/internal/cache"
)

type UseCase interface {
	GetMusicChart() (models.GetMusicChartResponse, error)
	GetRecommendedSongs(data models.GetRecommendedSongsRequest) (models.GetRecommendedSongsResponse, error)
	InsertAction(data models.ActionMessage) error
}

type Usecase struct {
	cfg   config.Config
	Repo  repository.Repository
	Cache cache.Cache
}

func NewUsecase(database *sql.DB, cache *redis.Client, cfg config.Config) UseCase {
	pgRepo := repository.NewPGRepository(database)
	redisChache := rediscashe.NewRedisCache(cache)

	return &Usecase{Repo: pgRepo, Cache: redisChache, cfg: cfg}
}

func (u *Usecase) GetMusicChart() (models.GetMusicChartResponse, error) {
	var res models.GetMusicChartResponse
	var err error

	data, cacheErr := u.Cache.Get(string(models.Chart))
	if cacheErr == nil {
		err = json.Unmarshal(data, &res)
		if err == nil {
			logrus.Info("Return from cache")
			return res, nil
		} else {
			logrus.Error(err)
		}
	} else if !errors.Is(cacheErr, redis.Nil) {
		logrus.Error(cacheErr)
	}

	res, err = u.Repo.GetMusicChart()
	if err != nil {
		return res, nil
	}

	bytes, cacheErr := json.Marshal(res)
	cacheErr = u.Cache.Set(string(models.Chart), bytes, u.cfg.Redis.SaveTimeHour)
	if cacheErr != nil {
		logrus.Error(cacheErr)
	}

	return res, nil
}

func (u *Usecase) GetRecommendedSongs(req models.GetRecommendedSongsRequest) (models.GetRecommendedSongsResponse, error) {
	var res models.GetRecommendedSongsResponse
	var err error

	data, cacheErr := u.Cache.Get(string(models.Recs))
	if cacheErr == nil {
		err = json.Unmarshal(data, &res)
		if err == nil {
			logrus.Info("Return from cache")
			return res, nil
		} else {
			logrus.Error(err)
		}
	} else if !errors.Is(cacheErr, redis.Nil) {
		logrus.Error(cacheErr)
	}

	res, err = u.Repo.GetRecommendedSongs(req)
	if err != nil {
		return res, nil
	}

	bytes, cacheErr := json.Marshal(res)
	cacheErr = u.Cache.Set(string(models.Recs), bytes, u.cfg.Redis.SaveTimeHour)
	if cacheErr != nil {
		logrus.Error(cacheErr)
	}

	return res, nil
}

func (u *Usecase) InsertAction(data models.ActionMessage) error {
	err := u.Repo.InsertAction(data)
	if err != nil {
		return err
	}
	logrus.Infof("successful set data in postgress: UserID: %d, SongID %d, Time: %v, Action: %s",
		data.UserID, data.SongID, data.Time, data.Action)
	return nil
}

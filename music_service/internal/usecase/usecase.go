package usecase

import (
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/RecommendationService/config"
	"github.com/SerzhLimon/RecommendationService/internal/broker"
	"github.com/SerzhLimon/RecommendationService/internal/cache"
	"github.com/SerzhLimon/RecommendationService/internal/models"
	"github.com/SerzhLimon/RecommendationService/internal/repository"
)

type UseCase interface {
	LikeSong(data models.LikeSongRequest) error
	ListenSong(data models.ListenSongRequest) (models.ListenSongResponse, error)
	DeleteSong(data models.DeleteSongRequest) error
	CreateSong(data models.CreateSongRequest) (models.CreateSongResponse, error)
	UpdateSong(data models.UpdateSongRequest) error
}

type Usecase struct {
	cfg    config.Config
	Repo   repository.Repository
	Cache  cache.Cache
	Broker broker.Broker
}

func NewUsecase(repo repository.Repository, cache cache.Cache, broker broker.Broker, cfg config.Config) UseCase {
	return &Usecase{Repo: repo, Cache: cache, Broker: broker, cfg: cfg}
}

func (u *Usecase) ListenSong(data models.ListenSongRequest) (models.ListenSongResponse, error) {
	var err error
	defer func() {
		if errors.Is(err, redis.Nil) || err == nil {
			err := u.Broker.SendActionMessage(u.cfg.Kafka.TopicSong, data.UserID, data.SongID, models.Listen)
			if err != nil {
				logrus.Error(err)
			}
		}
	}()

	song, err := u.Cache.GetSong(data)
	if err == nil {
		logrus.Info("Return song from cache")
		return models.ListenSongResponse{Song: song}, nil
	}
	if !errors.Is(err, redis.Nil) {
		logrus.Error(err)
	}

	res, err := u.Repo.ListenSong(data)
	if err != nil {
		return res, err
	}

	cacheErr := u.Cache.SetSong(data, res.Song, u.cfg.Redis.SaveTimeHour)
	if cacheErr != nil {
		logrus.Error(cacheErr)
	}

	return res, err
}

func (u *Usecase) DeleteSong(data models.DeleteSongRequest) error {
	err := u.Repo.DeleteSong(data)
	if err != nil {
		return err
	}
	err = u.Broker.SendEventMessage(u.cfg.Kafka.TopicAnalitic, data.SongID, models.Delete)
	if err != nil {
		logrus.Error(err)
	}

	return nil
}

func (u *Usecase) CreateSong(data models.CreateSongRequest) (models.CreateSongResponse, error) {
	res, err := u.Repo.CreateSong(data)
	if err != nil {
		return res, err
	}

	err = u.Broker.SendEventMessage(u.cfg.Kafka.TopicAnalitic, res.SongID, models.Create)
	if err != nil {
		logrus.Error(err)
	}

	return res, nil
}

func (u *Usecase) UpdateSong(data models.UpdateSongRequest) error {
	err := u.Repo.UpdateSong(data)
	if err != nil {
		return err
	}

	err = u.Broker.SendEventMessage(u.cfg.Kafka.TopicAnalitic, data.SongID, models.Update)
	if err != nil {
		logrus.Error(err)
	}

	return nil
}

func (u *Usecase) LikeSong(data models.LikeSongRequest) error {
	return u.Broker.SendActionMessage(u.cfg.Kafka.TopicSong, data.UserID, data.SongID, models.Like)
}

func SafeDereference(str *string) string {
	if str == nil {
		return "nil"
	}
	return *str
}

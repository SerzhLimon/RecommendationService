package usecase

import (
	"database/sql"

	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/RecommendationService/config"
	"github.com/SerzhLimon/RecommendationService/internal/models"
	"github.com/SerzhLimon/RecommendationService/internal/repository"
)

type UseCase interface {
	InsertEventUser(data models.UserEventMessage) error
	InsertEventSong(data models.SongEventMessage) error
}

type Usecase struct {
	cfg  config.Config
	Repo repository.Repository
}

func NewUsecase(database *sql.DB, cfg config.Config) UseCase {
	pgRepo := repository.NewPGRepository(database)

	return &Usecase{Repo: pgRepo, cfg: cfg}
}

func (u *Usecase) InsertEventUser(data models.UserEventMessage) error {
	err := u.Repo.InsertEventUser(data)
	if err != nil {
		return err
	}
	logrus.Infof("successful set data in postgress: UserID: %d, Time: %v, Event: %s",
		data.UserID, data.Time, data.Event)
	return nil
}

func (u *Usecase) InsertEventSong(data models.SongEventMessage) error {
	err := u.Repo.InsertEventSong(data)
	if err != nil {
		return err
	}
	logrus.Infof("successful set data in postgress: SongID: %d, Time: %v, Event: %s",
		data.SongID, data.Time, data.Event)
	return nil
}

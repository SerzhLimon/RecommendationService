package usecase

import (
	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/RecommendationService/config"
	"github.com/SerzhLimon/RecommendationService/internal/broker"
	"github.com/SerzhLimon/RecommendationService/internal/models"
	"github.com/SerzhLimon/RecommendationService/internal/repository"
)

type UseCase interface {
	CreateUser(data models.CreateUserRequest) (models.CreateUserResponse, error)
	UpdateUser(data models.UpdateUserRequest) error
	DeleteUser(data models.DeleteUserRequest) error
	GetUser(data models.GetUserRequest) (models.GetUserResponse, error)
}

type Usecase struct {
	cfg    config.Config
	Repo   repository.Repository
	Broker broker.Broker
}

func NewUsecase(repo repository.Repository, broker broker.Broker, cfg config.Config) UseCase {
	return &Usecase{Repo: repo, Broker: broker, cfg: cfg}
}

func SafeDereference(str *string) string {
	if str == nil {
		return "nil"
	}
	return *str
}

func (u *Usecase) CreateUser(data models.CreateUserRequest) (models.CreateUserResponse, error) {
	res, err := u.Repo.CreateUser(data)
	if err != nil {
		return res, err
	}

	err = u.Broker.SendMessage(u.cfg.Kafka.Topic, res.UserID, models.Create)
	if err != nil {
		logrus.Error(err)
	}

	return res, nil
}

func (u *Usecase) UpdateUser(data models.UpdateUserRequest) error {
	err := u.Repo.UpdateUser(data)
	if err != nil {
		return err
	}

	err = u.Broker.SendMessage(u.cfg.Kafka.Topic, data.UserID, models.Update)
	if err != nil {
		logrus.Error(err)
	}

	return nil
}

func (u *Usecase) DeleteUser(data models.DeleteUserRequest) error {
	err := u.Repo.DeleteUser(data)
	if err != nil {
		return err
	}

	err = u.Broker.SendMessage(u.cfg.Kafka.Topic, data.UserID, models.Delete)
	if err != nil {
		logrus.Error(err)
	}

	return nil
}

func (u *Usecase) GetUser(data models.GetUserRequest) (models.GetUserResponse, error) {
	return u.Repo.GetUser(data)
}

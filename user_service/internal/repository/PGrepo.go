package repository

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/SerzhLimon/RecommendationService/internal/models"
)

type Repository interface {
	CreateUser(data models.CreateUserRequest) (models.CreateUserResponse, error)
	UpdateUser(data models.UpdateUserRequest) error
	DeleteUser(data models.DeleteUserRequest) error
	GetUser(data models.GetUserRequest) (models.GetUserResponse, error)
}

type pgRepo struct {
	db *sql.DB
}

func NewPGRepository(db *sql.DB) Repository {
	return &pgRepo{db: db}
}

func (r *pgRepo) CreateUser(data models.CreateUserRequest) (models.CreateUserResponse, error) {
	var userID int64
	if err := r.db.QueryRow(queryCreateUser, data.UserName, data.PasswordHash).Scan(&userID); err != nil {
		err := errors.Errorf("pgRepo.CreateUser %v", err)
		return models.CreateUserResponse{}, err
	}

	return models.CreateUserResponse{UserID: userID}, nil
}

func (r *pgRepo) UpdateUser(data models.UpdateUserRequest) error {
	result, err := r.db.Exec(queryUpdateUser,
		data.UserID,
		data.NewName,
		data.NewPasswordHash,
	)
	if err != nil {
		err := errors.Errorf("pgRepo.UpdateUser %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Errorf("pgRepo.UpdateUser: %v", err)
	}

	if rowsAffected < 1 {
		err = errors.New("pgRepo.UpdateUser: no rows updated")
		return err
	}

	return nil
}

func (r *pgRepo) DeleteUser(data models.DeleteUserRequest) error {
	res, err := r.db.Exec(queryDeleteUser, data.UserID)
	if err != nil {
		err := errors.Errorf("pgRepo.DeleteUser %v", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		err := errors.Errorf("pgRepo.DeleteUser %v", err)
		return err
	}
	if rowsAffected < 1 {
		err := errors.Errorf("pgRepo.DeleteUser: no rows affected in users, user not found")
		return err
	}

	return nil
}

func (r *pgRepo) GetUser(data models.GetUserRequest) (models.GetUserResponse, error) {
	var res models.GetUserResponse

	err := r.db.QueryRow(queryGetUser, data.UserID).Scan(
		&res.UserName,
		&res.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, err
		}
		err = errors.Errorf("pgRepo.GetUser %v", err)
		return res, err
	}

	return res, nil
}

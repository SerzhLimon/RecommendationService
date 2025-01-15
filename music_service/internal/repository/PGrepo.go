package repository

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/RecommendationService/internal/models"
)

type Repository interface {
	ListenSong(data models.ListenSongRequest) (models.ListenSongResponse, error)
	DeleteSong(data models.DeleteSongRequest) error
	CreateSong(data models.CreateSongRequest) (models.CreateSongResponse, error)
	UpdateSong(data models.UpdateSongRequest) error
}

type pgRepo struct {
	db *sql.DB
}

func NewPGRepository(db *sql.DB) Repository {
	return &pgRepo{db: db}
}

func (r *pgRepo) ListenSong(data models.ListenSongRequest) (models.ListenSongResponse, error) {

	var res models.ListenSongResponse
	err := r.db.QueryRow(queryListenSong, data.SongID).Scan(&res.Song)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, err
		}
		err = errors.Errorf("pgRepo.ListenSong %v", err)
		return res, err
	}

	return res, nil
}

func (r *pgRepo) DeleteSong(data models.DeleteSongRequest) error {
	res, err := r.db.Exec(queryDeleteSong, data.SongID)
	if err != nil {
		err := errors.Errorf("pgRepo.DeleteSong %v", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		err := errors.Errorf("pgRepo.DeleteSong %v", err)
		return err
	}
	if rowsAffected < 1 {
		err := errors.Errorf("pgRepo.DeleteSong: no rows affected in songs, song not found")
		return err
	}
	logrus.Infof("Rows affected: %d", rowsAffected)

	return nil
}

func (r *pgRepo) CreateSong(data models.CreateSongRequest) (models.CreateSongResponse, error) {
	var songID int64
	if err := r.db.QueryRow(queryCreateSong, data.SongName, data.Content).Scan(&songID); err != nil {
		err := errors.Errorf("pgRepo.CreateSong %v", err)
		return models.CreateSongResponse{}, err
	}

	return models.CreateSongResponse{SongID: songID}, nil
}

func (r *pgRepo) UpdateSong(data models.UpdateSongRequest) error {

	result, err := r.db.Exec(queryUpdateSong,
		data.SongID,
		data.NewName,
		data.NewContent,
	)
	if err != nil {
		err := errors.Errorf("pgRepo.UpdateSong %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Errorf("pgRepo.UpdateSong: %v", err)
	}

	if rowsAffected < 1 {
		err = errors.New("pgRepo.UpdateSong: no rows updated")
		return err
	}

	return nil
}

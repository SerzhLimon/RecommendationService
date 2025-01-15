package repository

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/SerzhLimon/RecommendationService/internal/models"
)

type Repository interface {
	InsertEventUser(data models.UserEventMessage) error
	InsertEventSong(data models.SongEventMessage) error
}

type pgRepo struct {
	db *sql.DB
}

func NewPGRepository(db *sql.DB) Repository {
	return &pgRepo{db: db}
}

func (r *pgRepo) InsertEventUser(data models.UserEventMessage) error {
	res, err := r.db.Exec(queryInsertEventUser, data.UserID, data.Time, data.Event)
	if err != nil {
		err := errors.Errorf("pgRepo.InsertEventUser %v", err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		err := errors.Errorf("pgRepo.InsertEventUser %v", err)
		return err
	}
	if rowsAffected < 1 {
		err := errors.Errorf("pgRepo.InsertEventUser: no rows affected in user_events")
		return err
	}

	return nil
}

func (r *pgRepo) InsertEventSong(data models.SongEventMessage) error {
	res, err := r.db.Exec(queryInsertEventSong, data.SongID, data.Time, data.Event)
	if err != nil {
		err := errors.Errorf("pgRepo.InsertEventSong %v", err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		err := errors.Errorf("pgRepo.InsertEventSong %v", err)
		return err
	}
	if rowsAffected < 1 {
		err := errors.Errorf("pgRepo.InsertEventSong: no rows affected in music_events")
		return err
	}

	return nil
}

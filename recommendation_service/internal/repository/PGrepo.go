package repository

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/SerzhLimon/RecommendationService/internal/models"
)

type Repository interface {
	GetMusicChart() (models.GetMusicChartResponse, error)
	GetRecommendedSongs(data models.GetRecommendedSongsRequest) (models.GetRecommendedSongsResponse, error)
	InsertAction(data models.ActionMessage) error
}

type pgRepo struct {
	db *sql.DB
}

func NewPGRepository(db *sql.DB) Repository {
	return &pgRepo{db: db}
}

func (r *pgRepo) GetMusicChart() (models.GetMusicChartResponse, error) {
	var res models.GetMusicChartResponse
	rows, err := r.db.Query(queryGetMusicChart)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, nil
		}
		err := errors.Errorf("pgRepo.GetMusicChart %v", err)
		return res, err
	}

	for rows.Next() {
		var song models.Song
		err := rows.Scan(
			&song.SongID,
			&song.SongName,
		)
		if err != nil {
			err := errors.Errorf("pgRepo.GetMusicChart %v", err)
			return res, err
		}
		res.Songs = append(res.Songs, song)
	}

	if err := rows.Err(); err != nil {
		err = errors.Errorf("pgRepo.GetMusicChart %v", err)
		return res, err
	}

	return res, nil
}

func (r *pgRepo) GetRecommendedSongs(data models.GetRecommendedSongsRequest) (models.GetRecommendedSongsResponse, error) {
	var res models.GetRecommendedSongsResponse
	rows, err := r.db.Query(queryGetRecommendedSongs, data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, nil
		}
		err := errors.Errorf("pgRepo.GetRecommendedSongs %v", err)
		return res, err
	}

	for rows.Next() {
		var song models.Song
		err := rows.Scan(
			&song.SongID,
			&song.SongName,
		)
		if err != nil {
			err := errors.Errorf("pgRepo.GetRecommendedSongs %v", err)
			return res, err
		}
		res.Songs = append(res.Songs, song)
	}

	if err := rows.Err(); err != nil {
		err = errors.Errorf("pgRepo.GetRecommendedSongs %v", err)
		return res, err
	}

	return res, nil
}

func (r *pgRepo) InsertAction(data models.ActionMessage) error {

	res, err := r.db.Exec(queryInsertAction, data.UserID, data.SongID, data.Time, data.Action)
	if err != nil {
		err := errors.Errorf("pgRepo.InsertAction %v", err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		err := errors.Errorf("pgRepo.InsertAction %v", err)
		return err
	}
	if rowsAffected < 1 {
		err := errors.Errorf("pgRepo.InsertAction: no rows affected in actions")
		return err
	}

	return nil
}

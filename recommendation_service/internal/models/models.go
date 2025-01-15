package models

import "time"

type GetMusicChartResponse struct {
	Songs []Song `json:"songs"`
}

type Song struct {
	SongID   int64  `json:"song_id"`
	SongName string `json:"song_name"`
}

type GetRecommendedSongsRequest struct {
	UserID int64 `form:"user_id" binding:"required"`
}

type GetRecommendedSongsResponse struct {
	Songs []Song `json:"songs"`
}

type Action string

const (
	Like   Action = "like"
	Listen Action = "listen"
)

type ActionMessage struct {
	UserID int64     `json:"user_id"`
	SongID int64     `json:"song_id"`
	Time   time.Time `json:"time"`
	Action Action    `json:"action"`
}

type Recommendation string

const (
	Chart Recommendation = "Chart"
	Recs  Recommendation = "Recs"
)

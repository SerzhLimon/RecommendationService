package models

import "time"

type LikeSongRequest struct {
	UserID int64 `json:"user_id"`
	SongID int64 `json:"song_id"`
}

type ListenSongRequest struct {
	UserID int64 `form:"user_id" binding:"required"`
	SongID int64 `form:"song_id" binding:"required"`
}

type ListenSongResponse struct {
	Song string `json:"song"`
}

type DeleteSongRequest struct {
	SongID int64 `json:"song_id"`
}

type CreateSongRequest struct {
	SongName string `json:"name"`
	Content  string `json:"content"`
}

type CreateSongResponse struct {
	SongID int64 `json:"song_id"`
}

type UpdateSongRequest struct {
	SongID     int64   `json:"song_id"`
	NewName    *string `json:"name,omitempty"`
	NewContent *string `json:"content,omitempty"`
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

type Event string

const (
	Create Event = "create"
	Update Event = "update"
	Delete Event = "delete"
)

type EventMessage struct {
	SongID int64     `json:"song_id"`
	Time   time.Time `json:"time"`
	Event  Event     `json:"event"`
}

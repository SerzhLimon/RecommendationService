package models

import "time"

type Event string

const (
	Create Event = "create"
	Update Event = "update"
	Delete Event = "delete"
)

type UserEventMessage struct {
	UserID int64     `json:"user_id"`
	Time   time.Time `json:"time"`
	Event  Event     `json:"event"`
}

type SongEventMessage struct {
	SongID int64     `json:"song_id"`
	Time   time.Time `json:"time"`
	Event  Event     `json:"event"`
}

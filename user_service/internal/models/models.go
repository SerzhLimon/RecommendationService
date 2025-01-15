package models

import "time"

type CreateUserRequest struct {
	UserName     string `json:"user_name"`
	PasswordHash string `json:"password_hash"`
}

type CreateUserResponse struct {
	UserID int64 `json:"user_id"`
}

type UpdateUserRequest struct {
	UserID          int64   `json:"user_id"`
	NewName         *string `json:"name,omitempty"`
	NewPasswordHash *string `json:"password_hash,omitempty"`
}

type DeleteUserRequest struct {
	UserID int64 `json:"user_id"`
}

type GetUserRequest struct {
	UserID int64 `form:"user_id" binding:"required"`
}

type GetUserResponse struct {
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}

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

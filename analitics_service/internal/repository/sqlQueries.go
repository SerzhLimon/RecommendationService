package repository

const (
	queryInsertEventSong = `
		INSERT INTO music_events (song_id, timestamp, event)
		VALUES ($1, $2, $3)
	`
	queryInsertEventUser = `
		INSERT INTO user_events (user_id, timestamp, event)
		VALUES ($1, $2, $3)
	`
)
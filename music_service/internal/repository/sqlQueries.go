package repository

const (
	queryListenSong = `
		SELECT content
		FROM songs
		WHERE id = $1;
	`

	queryDeleteSong = `
		DELETE FROM songs
		WHERE id = $1;
	`

	queryCreateSong = `
		INSERT INTO songs (name, content, created_at)
		VALUES($1, $2, NOW())
		RETURNING id
	`

	queryUpdateSong = `
		UPDATE songs
		SET 
			name = COALESCE($2, name),
			content = COALESCE($3, content),
			updated_at = NOW()
		WHERE id = $1;
	`
)
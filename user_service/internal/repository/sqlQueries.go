package repository

const (
	queryCreateUser = `
		INSERT INTO users (name, password_hash, created_at)
		VALUES($1, $2, NOW())
		RETURNING id
	`

	queryUpdateUser = `
		UPDATE users
		SET
			name = COALESCE($2, name),
			password_hash = COALESCE($3, password_hash),
			updated_at = NOW()
		WHERE id = $1;
	`

	queryDeleteUser = `
		DELETE FROM users
		WHERE id = $1;
	`

	queryGetUser = `
		SELECT name, created_at
		FROM users
		WHERE id = $1
	` 
)
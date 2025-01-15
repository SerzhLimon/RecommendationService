package migrations

const (
	create = `
		CREATE TABLE IF NOT EXISTS users
		(
			id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			name varchar(255) NOT NULL,
			password_hash VARCHAR(255) NOT NULL,                                                                                                             
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE
		);
	`
	drop = `
		DROP TABLE IF EXISTS users;
	`
	insertTestUsers = `
		INSERT INTO users (name, password_hash, created_at)
		VALUES
			('User 1', 'random symbols', NOW()),
			('User 2', 'random symbols', NOW()),
			('User 3', 'random symbols', NOW()),
			('User 4', 'random symbols', NOW()),
			('User 5', 'random symbols', NOW()),
			('User 6', 'random symbols', NOW()),
			('User 7', 'random symbols', NOW()),
			('User 8', 'random symbols', NOW()),
			('User 9', 'random symbols', NOW()),
			('User 10', 'random symbols', NOW());
	`
)
package migrations

const (
	createTable = `
		CREATE TABLE IF NOT EXISTS songs
		(
			id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			name varchar(255) NOT NULL,                                                                
			content text  NOT NULL,                                                
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE
		);
	`

	dropTable = `
		DROP TABLE IF EXISTS songs;
	`

	insertTestSongs = `
		INSERT INTO songs (name, content, created_at)
		VALUES
			('Song 1', 'Content of song 1', NOW()),
			('Song 2', 'Content of song 2', NOW()),
			('Song 3', 'Content of song 3', NOW()),
			('Song 4', 'Content of song 4', NOW()),
			('Song 5', 'Content of song 5', NOW()),
			('Song 6', 'Content of song 6', NOW()),
			('Song 7', 'Content of song 7', NOW()),
			('Song 8', 'Content of song 8', NOW()),
			('Song 9', 'Content of song 9', NOW()),
			('Song 10', 'Content of song 10', NOW());
	`
)
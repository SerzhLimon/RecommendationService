package migrations

const (
	createType = `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'events') THEN
				CREATE TYPE events AS ENUM ('create', 'update', 'delete');
			END IF;
		END $$;
	`
	createMusicEvents = `
		CREATE TABLE IF NOT EXISTS music_events
		(
			id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			song_id BIGINT NOT NULL,
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
			event events NOT NULL
		);
	`
	createUserEvents = `
		CREATE TABLE IF NOT EXISTS user_events
		(
			id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			user_id BIGINT NOT NULL,
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
			event events NOT NULL
		);
	`
	dropType = `
		DROP TYPE IF EXISTS events;		
	`
	dropUserEvents = `
		DROP TABLE IF EXISTS user_events;	
	`
	dropMusicEvents = `
		DROP TABLE IF EXISTS music_events;
	`
)
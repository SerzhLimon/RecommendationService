package migrations

const (
	createType = `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_action') THEN
				CREATE TYPE user_action AS ENUM ('listen', 'like');
			END IF;
		END $$;
	`
	createTable = `
		CREATE TABLE IF NOT EXISTS actions
		(
			id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			user_id BIGINT NOT NULL,
			song_id BIGINT NOT NULL,
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
			action user_action NOT NULL
		);
	`

	dropType = `
		DROP TYPE IF EXISTS user_action;	
	`
	dropTable = `
		DROP TABLE IF EXISTS actions;
	`

	insertTestActionListen = `
		INSERT INTO actions (user_id, song_id, timestamp, action)
		VALUES

			(1, 1, NOW() - INTERVAL '5 days', 'listen'),  
			(1, 2, NOW() - INTERVAL '10 days', 'listen'), 
			(2, 2, NOW() - INTERVAL '20 days', 'listen'), 
			(2, 3, NOW() - INTERVAL '15 days', 'listen'), 
			(3, 1, NOW() - INTERVAL '25 days', 'listen'),
			(3, 4, NOW() - INTERVAL '2 days', 'listen'),  
			(4, 5, NOW() - INTERVAL '8 days', 'listen'),  
			(5, 6, NOW() - INTERVAL '3 days', 'listen'),  
			(6, 7, NOW() - INTERVAL '7 days', 'listen'),  
			(7, 8, NOW() - INTERVAL '6 days', 'listen'),  
			(8, 9, NOW() - INTERVAL '9 days', 'listen'),  
			(9, 10, NOW() - INTERVAL '11 days', 'listen'),
			(10, 1, NOW() - INTERVAL '1 day', 'listen'),   
			
			(4, 1, NOW() - INTERVAL '4 days', 'listen'),  
			(5, 1, NOW() - INTERVAL '2 days', 'listen'),  
			(6, 1, NOW() - INTERVAL '1 days', 'listen'),  
			
			(1, 6, NOW() - INTERVAL '4 days', 'like'),    
			(2, 5, NOW() - INTERVAL '30 days', 'like'),   
		
			(10, 2, NOW() - INTERVAL '3 days', 'listen'), 
			(9, 3, NOW() - INTERVAL '5 days', 'listen'),  
			(8, 4, NOW() - INTERVAL '10 days', 'listen'), 
			(7, 5, NOW() - INTERVAL '15 days', 'listen'), 
		
			(1, 1, NOW() - INTERVAL '2 months', 'listen'); 
	`

	insertTestActionLike = `
		INSERT INTO actions (user_id, song_id, timestamp, action)
		VALUES
			-- likes by user
			(1, 1, NOW() - INTERVAL '5 days', 'like'),  
			(1, 2, NOW() - INTERVAL '10 days', 'like'),

			-- likes other users
			(2, 1, NOW() - INTERVAL '8 days', 'like'),  
			(2, 2, NOW() - INTERVAL '9 days', 'like'), 
			(3, 1, NOW() - INTERVAL '7 days', 'like'),  
			(4, 2, NOW() - INTERVAL '6 days', 'like'), 
			(5, 1, NOW() - INTERVAL '5 days', 'like'),  

			-- likes other users in other songs
			(2, 3, NOW() - INTERVAL '6 days', 'like'), 
			(3, 4, NOW() - INTERVAL '4 days', 'like'),  
			(4, 5, NOW() - INTERVAL '3 days', 'like'),  
			(4, 6, NOW() - INTERVAL '2 days', 'like'),  
			(5, 3, NOW() - INTERVAL '2 days', 'like'),  
			(6, 4, NOW() - INTERVAL '1 days', 'like'),  
			(7, 5, NOW() - INTERVAL '12 hours', 'like'), 
			(8, 6, NOW() - INTERVAL '8 hours', 'like'), 
			(2, 7, NOW() - INTERVAL '7 days', 'like'),  
			(3, 8, NOW() - INTERVAL '6 days', 'like'),  
			(4, 9, NOW() - INTERVAL '5 days', 'like'),  
			(5, 10, NOW() - INTERVAL '4 days', 'like'), 
			(6, 7, NOW() - INTERVAL '3 days', 'like'),  
			(7, 8, NOW() - INTERVAL '2 days', 'like'),  
			(8, 9, NOW() - INTERVAL '1 days', 'like'),  
			(9, 10, NOW() - INTERVAL '12 hours', 'like'), 

			-- other action
			(1, 3, NOW() - INTERVAL '5 days', 'listen'), 
			(2, 4, NOW() - INTERVAL '5 days', 'listen'), 
			(3, 5, NOW() - INTERVAL '5 days', 'listen'), 
			(4, 6, NOW() - INTERVAL '5 days', 'listen'), 

			-- old likes
			(5, 11, NOW() - INTERVAL '1 days', 'like'),  
			(6, 12, NOW() - INTERVAL '1 days', 'like'), 
			(7, 13, NOW() - INTERVAL '1 days', 'like'),  
			(8, 14, NOW() - INTERVAL '1 days', 'like'),  
			(2, 1, NOW() - INTERVAL '2 months', 'like'), 
			(3, 2, NOW() - INTERVAL '2 months', 'like'), 
			(4, 3, NOW() - INTERVAL '2 months', 'like'), 
			(5, 4, NOW() - INTERVAL '2 months', 'like'); 
	`
)
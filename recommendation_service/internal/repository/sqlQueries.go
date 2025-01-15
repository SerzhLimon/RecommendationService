package repository

const (

	//	top 100 songs with unique listen since one month

	queryGetMusicChart = `
		SELECT a.song_id, s.name AS song_name
		FROM actions a
		INNER JOIN songs s ON a.song_id = s.id
		WHERE a.action = 'listen' AND a.timestamp >= NOW() - INTERVAL '1 month'
		GROUP BY a.song_id, s.name
		LIMIT 100;
	`

	//	top 50 songs liked other users

	queryGetRecommendedSongs = `
		WITH UserLikedSongs AS (  --fetch songs liked user
				SELECT song_id
				FROM actions
				WHERE user_id = $1 AND action = 'like'
		),
		OtherUsers AS (          --fetch users who liked same songs
				SELECT DISTINCT user_id
				FROM actions
				WHERE 
						action = 'like' AND 
						song_id IN (SELECT song_id FROM UserLikedSongs) AND 
						user_id != $1
		),
		RecommendedSongs AS (   --fetch recommended songs
				SELECT song_id, COUNT(*) AS like_count
				FROM actions
				WHERE 
						action = 'like' AND 
						user_id IN (SELECT user_id FROM OtherUsers) AND 
						song_id NOT IN (SELECT song_id FROM UserLikedSongs)
				GROUP BY song_id
		)
		SELECT r.song_id, s.name AS song_name   --return top 50 rec songs
		FROM RecommendedSongs r
		INNER JOIN songs s ON r.song_id = s.id
		ORDER BY like_count DESC
		LIMIT 50;
	`

	queryInsertAction = `
		INSERT INTO actions (user_id, song_id, timestamp, action)
		VALUES ($1, $2, $3, $4)
	`

)
package transport

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/RecommendationService/config"
	"github.com/SerzhLimon/RecommendationService/internal/broker"
	"github.com/SerzhLimon/RecommendationService/internal/models"
	"github.com/SerzhLimon/RecommendationService/internal/repository"
	uc "github.com/SerzhLimon/RecommendationService/internal/usecase"
	rediscashe "github.com/SerzhLimon/RecommendationService/internal/cache"
)

type Server struct {
	Usecase uc.UseCase
}

func NewServer(database *sql.DB, cache *redis.Client, producer sarama.SyncProducer, cfg config.Config) *Server {
	pgClient := repository.NewPGRepository(database)
	redisClient := rediscashe.NewRedisCache(cache)
	brokerClient := broker.NewKafkaClient(producer)
	uc := uc.NewUsecase(pgClient, redisClient, brokerClient, cfg)

	return &Server{Usecase: uc}
}

// ListenSong godoc
// @Summary      Listen to a song
// @Description  Simulates listening to a song by providing the user ID and song ID as query parameters. Returns the song's text if found.
// @Tags         music
// @Accept       json
// @Produce      json
// @Param        user_id query int64 true "User ID"
// @Param        song_id query int64 true "Song ID"
// @Success      200 {object} models.ListenSongResponse "Song details"
// @Failure      400 {object} map[string]string "error: invalid query parameters or incorrect IDs"
// @Failure      404 {object} map[string]string "error: song not found"
// @Router       /song [get]
func (s *Server) ListenSong(c *gin.Context) {
	var request models.ListenSongRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		logrus.WithError(err).Error("error binding query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	if request.SongID < 1 || request.UserID < 1 {
		logrus.Warn("validation failed: song's ID and user's ID must be greater than 0")
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect request"})
		return
	}

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("Parsed request: songID: %d, userID: %d", request.SongID, request.UserID)

	song, err := s.Usecase.ListenSong(request)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "songs not found"})
			return
		}
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "fail to listen song"})
		return
	}

	c.JSON(http.StatusOK, song)
}

// LikeSong godoc
// @Summary      Like a song
// @Description  Allows a user to like a song by providing the user ID and song ID in the request body.
// @Tags         music
// @Accept       json
// @Produce      json
// @Param        request body models.LikeSongRequest true "Like Song Request"
// @Success      200 {object} map[string]string "success: true"
// @Failure      400 {object} map[string]string "error: invalid request or failed to like song"
// @Router       /song/like [post]
func (s *Server) LikeSong(c *gin.Context) {
	var request models.LikeSongRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.WithError(err).Error("Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}
	if request.SongID < 1 || request.UserID < 1 {
		logrus.Warn("validation failed: song's ID and user's ID must be greater than 0")
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect request"})
		return
	}

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("Parsed request: %d", request.SongID)

	err := s.Usecase.LikeSong(request)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "fail to like song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "true"})
}

// UpdateSong godoc
// @Summary      Update an existing song
// @Description  This endpoint allows updating an existing song's details like name and/or text. You must provide the song's ID and the fields to be updated (name and/or text).
// @Tags         music
// @Accept       json
// @Produce      json
// @Param        request body models.UpdateSongRequest true "Song update request"
// @Success      200 {object} map[string]string "success: true"
// @Failure      400 {object} map[string]string "error: invalid request, incorrect song ID or failed to update song"
// @Router      /song/update [patch]
func (s *Server) UpdateSong(c *gin.Context) {
	var request models.UpdateSongRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.WithError(err).Error("error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format"})
		return
	}

	if request.SongID < 1 {
		logrus.Warn("validation failed: song's ID must be greater than 0")
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect song's ID"})
		return
	}

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("Parsed request: %d %s %s",
		request.SongID,
		uc.SafeDereference(request.NewName),
		uc.SafeDereference(request.NewContent),
	)

	err := s.Usecase.UpdateSong(request)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "fail to update song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "true"})
}

// DeleteSong godoc
// @Summary      Delete an existing song
// @Description  This endpoint allows deleting a song by providing their song ID. The song ID must be greater than 0 for a successful deletion.
// @Tags         music
// @Accept       json
// @Produce      json
// @Param        request body models.DeleteSongRequest true "Song deletion request"
// @Success      200 {object} map[string]string "success: true"
// @Failure      400 {object} map[string]string "error: invalid request, incorrect song ID or failed to delete song"
// @Router       /song/delete [delete]
func (s *Server) DeleteSong(c *gin.Context) {
	var request models.DeleteSongRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.WithError(err).Error("Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}
	if request.SongID < 1 {
		logrus.Warn("validation failed: song's ID must be greater than 0")
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect song's ID"})
		return
	}

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("Parsed request: %d", request.SongID)

	err := s.Usecase.DeleteSong(request)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "fail to delete song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "true"})
}

// CreateSong godoc
// @Summary      Create a new song
// @Description  This endpoint allows creating a new song with a songname and a songtext.
// @Tags         music
// @Accept       json
// @Produce      json
// @Param        request body models.CreateSongRequest true "Song creation request"
// @Success      200 {object} models.CreateSongResponse
// @Failure      400 {object} map[string]string "error: invalid request or failed to create song"
// @Router       /song/create [post]
func (s *Server) CreateSong(c *gin.Context) {
	var request models.CreateSongRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.WithError(err).Error("Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}
	if request.Content == "" || request.SongName == "" {
		logrus.Warn("validation failed: song's name or content is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect request"})
		return
	}

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("Parsed request: name: %s, content: %s", request.SongName, request.Content)

	res, err := s.Usecase.CreateSong(request)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "fail to create song"})
		return
	}

	c.JSON(http.StatusOK, res)
}

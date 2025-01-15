package transport

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/RecommendationService/config"
	"github.com/SerzhLimon/RecommendationService/internal/broker"
	"github.com/SerzhLimon/RecommendationService/internal/models"
	"github.com/SerzhLimon/RecommendationService/internal/repository"
	uc "github.com/SerzhLimon/RecommendationService/internal/usecase"
)

type Server struct {
	Usecase uc.UseCase
}

func NewServer(database *sql.DB, producer sarama.SyncProducer, cfg config.Config) *Server {
	pgClient := repository.NewPGRepository(database)
	brokerClient := broker.NewKafkaClient(producer)
	uc := uc.NewUsecase(pgClient, brokerClient, cfg)

	return &Server{Usecase: uc}
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  This endpoint allows creating a new user with a username and a hashed password.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body models.CreateUserRequest true "User creation request"
// @Success      200 {object} models.CreateUserResponse
// @Failure      400 {object} map[string]string "error: invalid request or failed to create user"
// @Router       /user/create [post]
func (s *Server) CreateUser(c *gin.Context) {
	var request models.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.WithError(err).Error("Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}
	if request.UserName == "" || request.PasswordHash == "" {
		logrus.Warn("validation failed: user's name or password hash is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect request"})
		return
	}

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("Parsed request: name: %s", request.UserName)

	res, err := s.Usecase.CreateUser(request)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "fail to create user"})
		return
	}

	c.JSON(http.StatusOK, res)
}

// UpdateUser godoc
// @Summary      Update an existing user
// @Description  This endpoint allows updating an existing user's details like name and/or password hash. You must provide the user's ID and the fields to be updated (name and/or password hash).
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body models.UpdateUserRequest true "User update request"
// @Success      200 {object} map[string]string "success: true"
// @Failure      400 {object} map[string]string "error: invalid request, incorrect user ID or failed to update user"
// @Router       /user/update [patch]
func (s *Server) UpdateUser(c *gin.Context) {
	var request models.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.WithError(err).Error("error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format"})
		return
	}

	if request.UserID < 1 {
		logrus.Warn("validation failed: users's ID must be greater than 0")
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect user's ID"})
		return
	}

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("Parsed request: %d %s %s",
		request.UserID,
		uc.SafeDereference(request.NewName),
		uc.SafeDereference(request.NewPasswordHash),
	)

	err := s.Usecase.UpdateUser(request)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "fail to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "true"})
}

// DeleteUser godoc
// @Summary      Delete an existing user
// @Description  This endpoint allows deleting a user by providing their user ID. The user ID must be greater than 0 for a successful deletion.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body models.DeleteUserRequest true "User deletion request"
// @Success      200 {object} map[string]string "success: true"
// @Failure      400 {object} map[string]string "error: invalid request, incorrect user ID or failed to delete user"
// @Router       /user/delete [delete]
func (s *Server) DeleteUser(c *gin.Context) {
	var request models.DeleteUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.WithError(err).Error("Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}
	if request.UserID < 1 {
		logrus.Warn("validation failed: user's ID must be greater than 0")
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect user's ID"})
		return
	}

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("Parsed request: %d", request.UserID)

	err := s.Usecase.DeleteUser(request)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "fail to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "true"})
}

// GetUser godoc
// @Summary      Retrieve user details
// @Description  This endpoint allows retrieving an existing user's details by providing their user ID as a query parameter. If the user does not exist, it returns a 404 error.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user_id query int64 true "User ID"
// @Success      200 {object} models.GetUserResponse "User details"
// @Failure      400 {object} map[string]string "error: invalid query parameters or incorrect user ID"
// @Failure      404 {object} map[string]string "error: user not found"
// @Router       /user [get]
func (s *Server) GetUser(c *gin.Context) {
	var request models.GetUserRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		logrus.WithError(err).Error("error binding query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	if request.UserID < 1 {
		logrus.Warn("validation failed: users's ID must be greater than 0")
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect user's ID"})
		return
	}

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("Parsed request: %d", request.UserID)

	res, err := s.Usecase.GetUser(request)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "fail to get user"})
		return
	}

	c.JSON(http.StatusOK, res)
}

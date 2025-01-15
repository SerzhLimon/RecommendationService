package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/RecommendationService/internal/models"
	uc "github.com/SerzhLimon/RecommendationService/internal/usecase"
)

type Server struct {
	Usecase uc.UseCase
}

func NewServer(uc uc.UseCase) *Server {
	return &Server{Usecase: uc}
}

// GetMusicChart godoc
// @Summary      Retrieve Music Chart
// @Description  This endpoint returns a list of songs currently in the music chart.
// @Tags         recommendations
// @Accept       json
// @Produce      json
// @Success      200 {object} models.GetMusicChartResponse "List of songs in the music chart"
// @Failure      400 {object} map[string]string "error: fail to get music chart"
// @Router       /chart [get]
func (s *Server) GetMusicChart(c *gin.Context) {
	chart, err := s.Usecase.GetMusicChart()
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "fail to get music chart"})
		return
	}
	c.JSON(http.StatusOK, chart)
}

// GetRecommendedSongs godoc
// @Summary      Get Recommended Songs
// @Description  This endpoint returns a list of recommended songs for a user based on their user ID.
// @Tags         recommendations
// @Accept       json
// @Produce      json
// @Param        user_id query int64 true "User ID"
// @Success      200 {object} models.GetRecommendedSongsResponse "List of recommended songs"
// @Failure      400 {object} map[string]string "error: invalid query parameters or failed to get recommended songs"
// @Router       /recommended [get]
func (s *Server) GetRecommendedSongs(c *gin.Context) {
	var request models.GetRecommendedSongsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		logrus.WithError(err).Error("error binding query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}
	if request.UserID < 1 {
		logrus.Warn("validation failed: song's ID must be greater than 0")
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect user's ID"})
		return
	}

	res, err := s.Usecase.GetRecommendedSongs(request)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "fail to get recommended songs"})
		return
	}

	c.JSON(http.StatusOK, res)
}

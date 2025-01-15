package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/RecommendationService/config"
	"github.com/SerzhLimon/RecommendationService/pkg/kafka"
	"github.com/SerzhLimon/RecommendationService/pkg/postgress"
	"github.com/SerzhLimon/RecommendationService/pkg/postgress/migrations"
	"github.com/SerzhLimon/RecommendationService/pkg/redis"
	serv "github.com/SerzhLimon/RecommendationService/internal/transport"
)

//	@title			Music Service
//	@version		1.0
//	@description	This is a simple music service.
//	@termsOfService	http://swagger.io/terms/

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8000
// @BasePath	/
// @schemes       http

func main() {

	gin.SetMode(gin.ReleaseMode)

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)

	logrus.Info("Loading configuration...")
	config := config.LoadConfig()
	logrus.Debugf("Configuration loaded: %+v", config)

	logrus.Info("Initializing PostgreSQL client...")
	db, err := postgress.InitPostgresClient(config)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize PostgreSQL client")
	}
	defer func() {
		logrus.Info("Closing PostgreSQL connection...")
		db.Close()
		logrus.Info("PostgreSQL connection closed")
	}()
	logrus.Info("PostgreSQL client initialized successfully")

	logrus.Info("Initializing Redis client...")
	redis, err := redis.InitRedisClient(config)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize Redis client")
	}
	defer func() {
		logrus.Info("Closing Redis connection...")
		redis.Close()
		logrus.Info("Redis connection closed")
	}()
	logrus.Info("Redis client initialized successfully")

	logrus.Info("Initializing Kafka client...")
	kafka, err := kafka.InitKafkaProducer(config)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize Kafka client")
	}
	logrus.Info("Kafka client initialized successfully")

	logrus.Info("Initializing server...")
	server := serv.NewServer(db, redis, kafka, config)
	routes := serv.ApiHandleFunctions{
		Server: *server,
	}

	logrus.Info("Running migrations...")
	err = migrations.Up(db)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to apply migrations")
	}
	defer func() {
		migrations.Down(db)
		logrus.Info("Migrations down")
	}()
	logrus.Info("Migrations applied successfully")

	logrus.Info("Setting up router...")
	router := serv.NewRouter(routes)

	logrus.Infof("Starting server on port %s...", ":8000")
	go func() {
		if err := router.Run(":8000"); err != nil {
			logrus.WithError(err).Fatal("Failed to start server")
		}
	}()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	<-sigchan
}

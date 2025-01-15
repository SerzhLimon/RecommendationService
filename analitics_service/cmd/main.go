package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/RecommendationService/config"
	"github.com/SerzhLimon/RecommendationService/internal/consumer"
	uc "github.com/SerzhLimon/RecommendationService/internal/usecase"
	"github.com/SerzhLimon/RecommendationService/pkg/kafka"
	"github.com/SerzhLimon/RecommendationService/pkg/postgress"
	"github.com/SerzhLimon/RecommendationService/pkg/postgress/migrations"
)

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

	uc := uc.NewUsecase(db, config)

	logrus.Info("Initializing Kafka consumer...")
	consumerGroup, err := kafka.InitKafkaConsumerGroup(config)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize Kafka consumer")
	}
	consumer := consumer.NewKafkaConsumer(uc, consumerGroup, config)

	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.StartKafkaConsumer(ctx)
	}()

	logrus.Infof("Starting server on port %s...", ":8003")
	server := gin.Default()
	go func() {
		if err := server.Run(":8003"); err != nil {
			logrus.WithError(err).Fatal("Failed to start server")
		}
	}()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	<-sigchan
	cancel()
	wg.Wait()
}

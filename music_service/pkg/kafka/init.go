package kafka

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/RecommendationService/config"
)

func InitKafkaProducer(cfg config.Config) (sarama.SyncProducer, error) {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = cfg.Kafka.MaxRetries
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.ClientID = cfg.Kafka.ClientID
	saramaConfig.Metadata.Retry.Backoff = time.Duration(cfg.Kafka.RetryBackoffMs) * time.Millisecond

	producer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, saramaConfig)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"brokers": cfg.Kafka.Brokers,
			"error":   err.Error(),
		}).Error("Failed to initialize Kafka producer")
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"brokers":               cfg.Kafka.Brokers,
		"topic for recomendate": cfg.Kafka.TopicSong,
		"topic for analitics":   cfg.Kafka.TopicAnalitic,
	}).Info("Successfully connected to Kafka")

	return producer, nil
}

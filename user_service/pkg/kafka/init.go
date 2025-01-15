package kafka

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/RecommendationService/config"
)

func InitKafkaProducer(cfg config.KafkaConfig) (sarama.SyncProducer, error) {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = cfg.MaxRetries
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.ClientID = cfg.ClientID
	saramaConfig.Metadata.Retry.Backoff = time.Duration(cfg.RetryBackoffMs) * time.Millisecond

	producer, err := sarama.NewSyncProducer(cfg.Brokers, saramaConfig)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"brokers": cfg.Brokers,
			"error":   err.Error(),
		}).Error("Failed to initialize Kafka producer")
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"brokers": cfg.Brokers,
		"topic":   cfg.Topic,
	}).Info("Successfully connected to Kafka")

	return producer, nil
}

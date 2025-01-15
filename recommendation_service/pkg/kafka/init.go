package kafka

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/RecommendationService/config"
)

func InitKafkaConsumer(cfg config.Config) (sarama.Consumer, error) {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	saramaConfig.Consumer.Group.Session.Timeout = time.Duration(cfg.Kafka.RetryBackoffMs) * time.Millisecond
	saramaConfig.Consumer.Group.Heartbeat.Interval = 10 * time.Millisecond
	saramaConfig.ClientID = cfg.Kafka.ClientID

	consumer, err := sarama.NewConsumer(cfg.Kafka.Brokers, saramaConfig)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"brokers": cfg.Kafka.Brokers,
			"error":   err.Error(),
		}).Error("Failed to initialize Kafka consumer")
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"brokers": cfg.Kafka.Brokers,
		"topic":   cfg.Kafka.Topic,
	}).Info("Successfully connected to Kafka")

	return consumer, nil
}

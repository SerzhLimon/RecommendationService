package kafka

import (
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/RecommendationService/config"
)

func InitKafkaConsumerGroup(cfg config.Config) (sarama.ConsumerGroup, error) {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Version = sarama.V2_1_0_0

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, cfg.Kafka.GroupID, config)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"brokers": cfg.Kafka.Brokers,
			"error":   err.Error(),
		}).Error("failed to initialize Kafka consumer")
		return nil, err
	}

	logrus.Info("Successfully connected to Kafka")

	return consumerGroup, nil
}

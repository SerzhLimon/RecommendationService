package consumer

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/RecommendationService/config"
	"github.com/SerzhLimon/RecommendationService/internal/models"
	uc "github.com/SerzhLimon/RecommendationService/internal/usecase"
)

type topic int

const (
	songs_to_analitics topic = 0
	users_to_analitics topic = 1
)

type Consumer interface {
	StartKafkaConsumer(ctx context.Context, topic string)
}

type KafkaConsumer struct {
	Usecase       uc.UseCase
	ConsumerGroup sarama.ConsumerGroup
	cfg           config.Config
}

func NewKafkaConsumer(usecase uc.UseCase, cg sarama.ConsumerGroup, cfg config.Config) *KafkaConsumer {
	return &KafkaConsumer{
		Usecase:       usecase,
		ConsumerGroup: cg,
		cfg:           cfg,
	}
}

func (k *KafkaConsumer) StartKafkaConsumer(ctx context.Context) error {
	err := k.ConsumerGroup.Consume(ctx, k.cfg.Kafka.Topics, k)
	return err
}

func (k *KafkaConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (k *KafkaConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	logrus.Info("Consumer is closed")
	return nil
}

func (k *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		switch claim.Topic() {
		case k.cfg.Kafka.Topics[songs_to_analitics]:
			if err := k.processSongEvents(message.Value); err != nil {
				logrus.Error(err)
			}
		case k.cfg.Kafka.Topics[users_to_analitics]:
			if err := k.processUserEvents(message.Value); err != nil {
				logrus.Error(err)
			}
		}
	}
	return nil
}

func (k *KafkaConsumer) processSongEvents(message []byte) error {
	var songEvent models.SongEventMessage
	if err := json.Unmarshal(message, &songEvent); err != nil {
		return errors.Errorf("KafkaConsumer.processSongEvents %v", err)
	}
	if err := k.Usecase.InsertEventSong(songEvent); err != nil {
		return err
	}
	return nil
}

func (k *KafkaConsumer) processUserEvents(message []byte) error {
	var userEvent models.UserEventMessage
	if err := json.Unmarshal(message, &userEvent); err != nil {
		return errors.Errorf("KafkaConsumer.processUserEvents %v", err)
	}
	if err := k.Usecase.InsertEventUser(userEvent); err != nil {
		return err
	}
	return nil
}

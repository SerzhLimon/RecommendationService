package consumer

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/RecommendationService/internal/models"
	uc "github.com/SerzhLimon/RecommendationService/internal/usecase"
)

type Consumer interface {
	readMessages(ctx context.Context, topic string) error
	StartKafkaConsumer(ctx context.Context, topic string)
}

type KafkaConsumer struct {
	consumer sarama.Consumer
	Usecase  uc.UseCase
}

func NewKafkaClient(consumer sarama.Consumer, uc uc.UseCase) Consumer {
	return &KafkaConsumer{consumer: consumer, Usecase: uc}
}

func (k *KafkaConsumer) readMessages(ctx context.Context, topic string) error {
	partitionConsumer, err := k.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return errors.Errorf("KafkaConsumer.ReadMessages %v", err)
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg, ok := <-partitionConsumer.Messages():
			if !ok {
				return errors.New("message channel closed")
			}
			var mes models.ActionMessage
			err := json.Unmarshal(msg.Value, &mes)
			if err != nil {
				logrus.Error(errors.Errorf("KafkaConsumer.readMessages %v", err))
				continue
			}
			err = k.Usecase.InsertAction(mes)
			if err != nil {
				logrus.Error(err)
			}
		case <-ctx.Done():
			logrus.Info("Consumer finished successfully")
			return nil
		}
	}
}

func (k *KafkaConsumer) StartKafkaConsumer(ctx context.Context, topic string) {
	err := k.readMessages(ctx, topic)
	if err != nil {
		logrus.Error(err)
	}
}

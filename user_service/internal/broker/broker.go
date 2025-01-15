package broker

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"

	"github.com/SerzhLimon/RecommendationService/internal/models"
)

type Broker interface {
	SendMessage(topic string, userID int64, event models.Event) error
}

type KafkaBroker struct {
	producer sarama.SyncProducer
}

func NewKafkaClient(producer sarama.SyncProducer) Broker {
	return &KafkaBroker{producer: producer}
}

func (k *KafkaBroker) SendMessage(topic string, userID int64, event models.Event) error {
	msg := k.newActionMessage(userID, event)
	samaraMsg, err := k.buildMessage(msg, topic)
	if err != nil {
		err := errors.Errorf("KafkaBroker.SendMessage %v", err)
		return err
	}

	_, _, err = k.producer.SendMessage(samaraMsg)
	if err != nil {
		err := errors.Errorf("KafkaBroker.SendMessage %v", err)
		return err
	}

	return nil
}

func (k *KafkaBroker) newActionMessage(userID int64, event models.Event) *models.UserEventMessage {
	return &models.UserEventMessage{
		UserID: userID,
		Time:   time.Now(),
		Event:  event,
	}
}

func (k *KafkaBroker) buildMessage(message *models.UserEventMessage, topic string) (*sarama.ProducerMessage, error) {
	msg, err := json.Marshal(message)

	if err != nil {
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic: topic,
		// Key:   sarama.StringEncoder(message.Action),
		Value: sarama.ByteEncoder(msg),
	}, nil
}

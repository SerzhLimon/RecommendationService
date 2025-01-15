package broker

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"

	"github.com/SerzhLimon/RecommendationService/internal/models"
)

type Broker interface {
	SendActionMessage(topic string, userID, songID int64, action models.Action) error
	SendEventMessage(topic string, songID int64, event models.Event) error
}

type KafkaBroker struct {
	producer sarama.SyncProducer
}

func NewKafkaClient(producer sarama.SyncProducer) Broker {
	return &KafkaBroker{producer: producer}
}

func (k *KafkaBroker) SendActionMessage(topic string, userID, songID int64, action models.Action) error {
	msg := k.newActionMessage(userID, songID, action)
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

func (k *KafkaBroker) SendEventMessage(topic string, songID int64, event models.Event) error {
	msg := k.newEventMessage(songID, event)
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

func (k *KafkaBroker) newActionMessage(userID, songID int64, action models.Action) *models.ActionMessage {
	return &models.ActionMessage{
		UserID: userID,
		SongID: songID,
		Time:   time.Now(),
		Action: action,
	}
}

func (k *KafkaBroker) newEventMessage(songID int64, event models.Event) *models.EventMessage {
	return &models.EventMessage{
		SongID: songID,
		Time:   time.Now(),
		Event:  event,
	}
}

func (k *KafkaBroker) buildMessage(message interface{}, topic string) (*sarama.ProducerMessage, error) {
	msg, err := json.Marshal(message)

	if err != nil {
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic: topic,
		// Key:       sarama.StringEncoder(message.Action),
		Value:     sarama.ByteEncoder(msg),
		Partition: 0,
	}, nil
}

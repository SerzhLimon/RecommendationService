package config

import (
	"encoding/json"
	"os"

	"github.com/sirupsen/logrus"
)

const filepath = "config/config.json"

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
	Password string `json:"password"`
}

type RedisConfig struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	SaveTimeHour int    `json:"save_time"`
}

type KafkaConfig struct {
    Brokers        []string `json:"brokers"`
    TopicSong      string   `json:"topic_song"`
    TopicAnalitic  string   `json:"topic_analitic"`
    ClientID       string   `json:"client_id"`
    GroupID        string   `json:"group_id"`
    MaxRetries     int      `json:"max_retries"`
    RetryBackoffMs int      `json:"retry_backoff_ms"`
}

type Config struct {
	Postgres PostgresConfig `json:"postgres"`
	Redis    RedisConfig    `json:"redis"`
	Kafka    KafkaConfig    `json:"kafka"`
}

func LoadConfig() Config {

	var config Config
	data, err := os.ReadFile(filepath)
	if err != nil {
		logrus.WithError(err).Fatal("can not load config")
	}
	if err = json.Unmarshal(data, &config); err != nil {
		logrus.WithError(err).Fatal("can not load config")
	}
	logrus.Info("Success load config")

	return config
}

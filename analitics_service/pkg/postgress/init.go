package postgress

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/RecommendationService/config"
)

func InitPostgresClient(cfg config.Config) (*sql.DB, error) {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	options := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.DBName, cfg.Postgres.Password, cfg.Postgres.SSLMode)

	database, err := sql.Open("postgres", options)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"host":    cfg.Postgres.Host,
			"port":    cfg.Postgres.Port,
			"user":    cfg.Postgres.User,
			"dbname":  cfg.Postgres.DBName,
			"sslmode": cfg.Postgres.SSLMode,
			"error":   err.Error(),
		}).Error("Failed to open PostgreSQL connection")
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"host":    cfg.Postgres.Host,
			"port":    cfg.Postgres.Port,
			"user":    cfg.Postgres.User,
			"dbname":  cfg.Postgres.DBName,
			"sslmode": cfg.Postgres.SSLMode,
			"error":   err.Error(),
		}).Error("Failed to ping PostgreSQL database")
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"host":    cfg.Postgres.Host,
		"port":    cfg.Postgres.Port,
		"user":    cfg.Postgres.User,
		"dbname":  cfg.Postgres.DBName,
		"sslmode": cfg.Postgres.SSLMode,
	}).Info("Successful connection to PostgreSQL")

	return database, nil
}

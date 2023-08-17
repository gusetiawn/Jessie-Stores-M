package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func New() *sql.DB {
	host := viper.GetString("DB_HOST")
	port := viper.GetString("DB_PORT")
	name := viper.GetString("DB_NAME")
	user := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	logger := logrus.
		WithField("func", "postgres.new").
		WithField("host", host).
		WithField("name", name).
		WithField("user", user).
		WithField("password", password)

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, name)
	db, err := sql.Open("postgres", url)
	if err != nil {
		logger.WithError(err).Error("error connecting to database")
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		logger.WithError(err).Error("error connecting to database")
		log.Fatal(err)
	}

	logger.Info("success")

	return db
}

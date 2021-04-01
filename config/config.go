package config

import (
	"errors"
	"github.com/caarlos0/env"
)

type API struct {
	ApiKey string `env:"API_KEY"`
	ApiSecret string `env:"API_SECRET"`
	AccessToken string `env:"ACCESS_TOKEN"`
	AccessTokenSecret string `env:"ACCESS_TOKEN_SECRET"`
}

type Kafka struct {
	Topic string `env:"TOPIC"`
	Broker string `env:"BROKER"`
}


func InitApi() (*API, error) {
	apiConfig := API{}
	if err := env.Parse(&apiConfig); err != nil {
		return nil, errors.New("cloud not load api auth tokens")
	}

	return &apiConfig, nil
}

func InitKafka() (*Kafka, error) {
	kafkaConfig := Kafka{}
	if err := env.Parse(&kafkaConfig); err != nil {
		return nil, errors.New("cloud not load kafka environment variables")
	}

	return &kafkaConfig, nil
}
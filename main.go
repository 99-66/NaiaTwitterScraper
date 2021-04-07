package main

import (
	"encoding/json"
	"github.com/99-66/NaiaTwitterScraper/config"
	"github.com/99-66/NaiaTwitterScraper/controller/kafka"
	"github.com/99-66/NaiaTwitterScraper/models"
	"github.com/Shopify/sarama"
	"log"
)

func main() {
	// API Token 정보를 초기화한다
	apiConfig, err := config.InitApi()
	if err != nil {
		panic(err)
	}

	// Kafka 정보를 초기화한다
	kafkaConfig, err := config.InitKafka()
	if err != nil {
		panic(err)
	}

	// 트윗을 전달할 채널을 생성한다
	c := make(chan models.Tweet)

	// 트위터 스트림에서 트윗을 채널로 전송한다
	go kafka.GetStreamTweets(apiConfig, c)

	// Test
	//go TestGetStreamTweets(apiConfig, c)

	// 비동기 프로듀서를 생성한다
	p, err := kafka.NewAsyncProducer(kafkaConfig.Broker)
	if err != nil {
		panic(err)
	}
	defer p.Close()

	// 채널에서 트윗을 받아온다
	for val := range c {
		//Kafka로 전달하기 위해 트윗을 Json으로 마샬링한다
		valJson, err := json.Marshal(val)
		if err != nil {
			log.Panicf("tweets failed json marshaling. %vn", err)
		}
		// Kafka로 트윗을 전달한다
		msg := &sarama.ProducerMessage{
			Topic: kafkaConfig.Topic,
			Value: sarama.ByteEncoder(valJson),
		}
		log.Printf("send to %v\n", val.Id)
		p.Input() <- msg
	}
}

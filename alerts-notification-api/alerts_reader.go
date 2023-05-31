package main

import (
	"log"
	"os"

	"github.com/Shopify/sarama"
)

type AlertsReader struct {
	Alerts chan []byte
	open   chan bool
}

func NewAlertsReader() *AlertsReader {
	return &AlertsReader{
		Alerts: make(chan []byte),
		open:   make(chan bool),
	}
}

func getEnvOrDefault(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func (ar *AlertsReader) Listen() {
	brokerList := []string{getEnvOrDefault("KAFKA_BROKER", "kafka:9092")}
	log.Printf("Kafka broker list: %v\n", brokerList)

	topic := getEnvOrDefault("KAFKA_TOPIC", "alerts")
	log.Printf("Kafka topic: %v\n", topic)

	groupID := getEnvOrDefault("KAFKA_GROUP_ID", "alerts-reader")
	log.Printf("Kafka group ID: %v\n", groupID)

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		log.Fatalf("Error creating consumer: %v\n", err)
	}

	for {
		partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
		if err != nil {
			log.Fatalf("Error creating partition consumer: %v\n", err)
		}

		for message := range partitionConsumer.Messages() {
			ar.Alerts <- message.Value
		}
	}
}

func (ar *AlertsReader) Close() {
	ar.open <- false
}

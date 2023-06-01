package main

import (
	"context"
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

	consumer := Consumer{
		ready:        make(chan bool),
		AlertsReader: ar,
	}
	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(brokerList, groupID, config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v\n", err)
	}
	topics := []string{topic}

	for {
		log.Println("Creating consumer session")
		if err := client.Consume(ctx, topics, &consumer); err != nil {
			log.Fatalf("Error creating partition consumer: %v\n", err)
		}

		if ctx.Err() != nil {
			break
		}
	}
	log.Println("Closing consumer")
	cancel()
}

func (ar *AlertsReader) Close() {
	ar.open <- false
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready        chan bool
	AlertsReader *AlertsReader
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		log.Printf("Waiting for message")
		select {
		case message := <-claim.Messages():
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			consumer.AlertsReader.Alerts <- message.Value
			session.MarkMessage(message, "")
			continue

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

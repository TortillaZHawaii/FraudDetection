package main

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/Shopify/sarama"
)

func main() {
	// Kafka broker addresses
	brokerList := []string{"localhost:9094"}

	// Configuration options for the Kafka producer
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 5                    // Retry up to 5 times to produce the message

	// Create a new Kafka async producer
	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatal("Failed to create Kafka producer:", err)
	}
	defer producer.Close()

	// Topic to write messages to
	topic := "test"

	// Channel to receive Kafka delivery reports
	deliveryReports := producer.Errors()

	// Channel to receive OS signals for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Wait group to wait for delivery report goroutine
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Goroutine to handle delivery reports
	go func() {
		defer wg.Done()
		for err := range deliveryReports {
			log.Printf("Failed to deliver message to Kafka: %v\n", err)
		}
	}()

	// Example message to send
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder("Hello, Kafka!"),
	}

	// Send the message to the Kafka topic
	producer.Input() <- message

	// Wait for OS signal for graceful shutdown
	<-signals

	// Close the producer and wait for delivery report goroutine to finish
	producer.AsyncClose()
	wg.Wait()

	log.Println("Kafka producer shut down gracefully")
}

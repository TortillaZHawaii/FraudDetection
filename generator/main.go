package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

func getEnvOrDefault(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func main() {
	// Kafka broker addresses
	brokerList := []string{getEnvOrDefault("KAFKA_BROKER", "kafka:9092")}
	log.Printf("Kafka broker list: %v\n", brokerList)

	cardOwnersCountEnv := getEnvOrDefault("CARD_OWNERS_COUNT", "100")
	cardCountEnv := getEnvOrDefault("CARD_COUNT", "1000")
	seedEnv := getEnvOrDefault("SEED", "0")
	waitTimeMsEnv := getEnvOrDefault("WAIT_TIME_MS", "1000")

	// parse env vars
	cardOwnersCount, err := strconv.ParseInt(cardOwnersCountEnv, 10, 64)
	if err != nil {
		panic(err)
	}
	log.Printf("Card owners count: %d\n", cardOwnersCount)

	cardCount, err := strconv.ParseInt(cardCountEnv, 10, 64)
	if err != nil {
		panic(err)
	}
	log.Printf("Card count: %d\n", cardCount)

	waitTimeMs, err := strconv.ParseInt(waitTimeMsEnv, 10, 64)
	if err != nil {
		panic(err)
	}
	log.Printf("Wait time (ms): %d\n", waitTimeMs)

	seed, err := strconv.ParseInt(seedEnv, 10, 64)
	if err != nil {
		panic(err)
	}
	if seed != 0 {
		log.Printf("Seed: %d\n", seed)
	} else {
		log.Printf("Seed was not set or was set to zero, using rand/crypto")
	}

	// Configuration options for the Kafka producer
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 5                    // Retry up to 5 times to produce the message

	// Create a new Kafka async producer
	log.Println("Creating Kafka producer...")

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatal("Failed to create Kafka producer:", err)
	}
	defer producer.Close()
	log.Println("Kafka producer created successfully")

	// Topic to write messages to
	topic := getEnvOrDefault("KAFKA_TOPIC", "transactions")

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

	ts := NewTransactionSource(seed, cardOwnersCount, cardCount)

	// Send the message to the Kafka topic
	go func() {
		for {
			waitTimeNs := waitTimeMs * 1_000_000
			<-time.After(time.Duration(waitTimeNs))
			transaction := ts.GetTransaction()
			messageValue, err := transaction.ToJSON()

			if err != nil {
				panic(err)
			}

			// Example message to send
			message := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.StringEncoder(string(messageValue)),
			}

			producer.Input() <- message
			log.Printf("Message %s sent to Kafka topic: %s\n", string(messageValue), topic)
		}
	}()

	// Wait for OS signal for graceful shutdown
	<-signals

	// Close the producer and wait for delivery report goroutine to finish
	producer.AsyncClose()
	wg.Wait()

	log.Println("Kafka producer shut down gracefully")
}

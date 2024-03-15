package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	BootstrapServers string `json:"bootstrap.servers"`
	Acks             string `json:"acks"`
	KeySerializer    string `json:"key.serializer"`
	ValueSerializer  string `json:"value.serializer"`
}

type Server struct {
	writer *kafka.Writer
}

func NewServer(writer *kafka.Writer) *Server {
	return &Server{writer: writer}
}

func LoadKafkaConfig(configFile string) (*KafkaConfig, error) {
	var config KafkaConfig
	file, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func signMessage(originalMessage []byte) []byte {
	return append(originalMessage, []byte("-signed")...)
}

func (s *Server) writeToKafkaWithRetry(msg kafka.Message, maxRetries int, initialBackoff time.Duration) error {
	var err error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if err = s.writer.WriteMessages(context.Background(), msg); err == nil {
			return nil
		}
		time.Sleep(backoff)
		backoff *= 2
		backoff += time.Duration(rand.Intn(100)) * time.Millisecond // Add jitter
	}
	return err
}

func main() {
	config, err := LoadKafkaConfig("kafka-config.json")
	if err != nil {
		panic(err)
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{config.BootstrapServers},
		Topic:          "message-sign-request",
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})
	defer r.Close()

	signedMessageWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{config.BootstrapServers},
		Topic:   "message-sign-complete",
	})
	defer signedMessageWriter.Close()

	server := NewServer(signedMessageWriter)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			// TODO: handle gracefully
			fmt.Printf("Error reading message: %s\n", err)
			continue
		}
		fmt.Printf("message_630pm at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		signedMessage := signMessage(m.Value)

		message := kafka.Message{
			Key:   m.Key,
			Value: signedMessage,
		}

		if err := server.writeToKafkaWithRetry(message, 5, 500*time.Millisecond); err != nil {
			fmt.Printf("failed to write signed message to Kafka: %s\n", err)
			continue
		}
	}
}

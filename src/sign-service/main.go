package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	BootstrapServers string `json:"bootstrap.servers"`
	Acks             string `json:"acks"`
	KeySerializer    string `json:"key.serializer"`
	ValueSerializer  string `json:"value.serializer"`
}

type App struct {
	writer *kafka.Writer
}

func NewApp(writer *kafka.Writer) *App {
	return &App{writer: writer}
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

func (app *App) writeToKafka(key, signedMessage []byte) error {
	return app.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   key,
			Value: signedMessage,
		},
	)
}

func main() {
	config, err := LoadKafkaConfig("kafka-config.json")
	if err != nil {
		panic(err)
	}
	fmt.Println("BootstrapServers:", config.BootstrapServers) // Debugging line
	if config.BootstrapServers == "" {
		panic("BootstrapServers is empty")
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.BootstrapServers},
		Topic:   "message-sign-request",
		// GroupID:        "my-group",
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

	app := NewApp(signedMessageWriter)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			// TODO: handle gracefully
			fmt.Printf("Error reading message: %s\n", err)
			continue
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		signedMessage := signMessage(m.Value)

		if err := app.writeToKafka(m.Key, signedMessage); err != nil {
			// TODO: handle gracefully
			fmt.Printf("failed to write signed message to Kafka: %s\n", err)
			// fmt.Printf("Error reading message: %s\n", err)
			continue
		}
	}
}

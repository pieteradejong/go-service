package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	BootstrapServers string `json:"bootstrap.servers"`
	Acks             string `json:"acks"`
	KeySerializer    string `json:"key.serializer"`
	ValueSerializer  string `json:"value.serializer"`
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

func main() {
	fmt.Println("Starting producer...")
	config, err := LoadKafkaConfig("../config/kafka-config.json")
	if err != nil {
		panic(err)
	}

	brokers := strings.Split(config.BootstrapServers, ",")

	w := kafka.NewWriter(kafka.WriterConfig{
		// Brokers: []string{config.BootstrapServers},
		Brokers: brokers,
		Topic:   "testtopic1",
	})
	defer w.Close()

	err = w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte("Hello from Go producer!"),
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sent message")
}

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
	config, err := LoadKafkaConfig("../../kafka-config.json")
	if err != nil {
		panic(err)
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{config.BootstrapServers},
		Topic:          "testtopic1",
		GroupID:        "my-group",
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			panic(err)
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		if err := r.CommitMessages(context.Background(), m); err != nil {
			panic("failed to commit messages: " + err.Error())
		}
	}
}

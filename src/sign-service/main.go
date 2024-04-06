package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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

func handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	clientChan := make(chan []byte)
	register <- clientChan

	defer func() {
		unregister <- clientChan
	}()

	for msg := range clientChan {
		fmt.Fprintf(w, "data: %s\n\n", msg)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}

var (
	register   = make(chan chan []byte)
	unregister = make(chan chan []byte)
	broadcast  = make(chan []byte)
)

func manageClients() {
	clients := make(map[chan []byte]bool)

	for {
		select {
		case client := <-register:
			clients[client] = true
		case client := <-unregister:
			delete(clients, client)
			close(client)
		case msg := <-broadcast:
			for client := range clients {
				client <- msg
			}
		}
	}
}

func main() {
	go manageClients()
	http.HandleFunc("/events", handleSSE)
	log.Println("http server started on :8000")
	go http.ListenAndServe(":8000", nil)

	config, err := LoadKafkaConfig("kafka-config.json")
	if err != nil {
		panic(err)
	}

	emojiCountReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{config.BootstrapServers},
		Topic:          "reaction-emoji-counts",
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})
	defer emojiCountReader.Close()

	for {
		m, err := emojiCountReader.ReadMessage(context.Background())
		if err != nil {
			// TODO: handle gracefully
			fmt.Printf("Error reading message: %s\n", err)
			continue
		}
		fmt.Printf("\033[32memoji count at offset %d: %s = %s\n\033[0m", m.Offset, string(m.Key), string(m.Value))

	}
}

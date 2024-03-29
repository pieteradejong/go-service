package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type HealthCheckResponse struct {
	Status   string `json:"status"`
	Producer string `json:"producer"`
	Consumer string `json:"consumer"`
}

type SignRequest struct {
	Message string `json:"message"`
}

type KafkaConfig struct {
	BootstrapServers string `json:"bootstrap.servers"`
	Acks             string `json:"acks"`
	KeySerializer    string `json:"key.serializer"`
	ValueSerializer  string `json:"value.serializer"`
}

type Server struct {
	KafkaWriter *kafka.Writer
}

func NewServer(kafkaWriter *kafka.Writer) *Server {
	return &Server{
		KafkaWriter: kafkaWriter,
	}
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

func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthCheckResponse{
		Status: "Healthy",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (s *Server) writeToKafkaWithRetry(msg kafka.Message, maxRetries int, initialBackoff time.Duration) error {
	var err error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if err = s.KafkaWriter.WriteMessages(context.Background(), msg); err == nil {
			return nil
		}
		time.Sleep(backoff)
		backoff *= 2
		backoff += time.Duration(rand.Intn(100)) * time.Millisecond
	}
	return err
}

func (s *Server) reactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
		return
	}

	var req SignRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	message := kafka.Message{
		Key:   []byte("Key-A"),
		Value: []byte(req.Message),
	}

	if err := s.writeToKafkaWithRetry(message, 5, 500*time.Millisecond); err != nil {
		http.Error(w, "Error sending message to Kafka", http.StatusInternalServerError)
		return
	}

	fmt.Println("Sent message:", req.Message)

	response := map[string]string{"message": "Message accepted for processing"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonResponse)
}

func main() {
	config, err := LoadKafkaConfig("kafka-config.json")
	if err != nil {
		panic(err)
	}

	// tlsConfig, err := tlsconfig.SetupTLSConfig("server.crt", "server.key", "server.crt")
	// if err != nil {
	// 	log.Fatalf("Failed to setup TLS config: %v", err)
	// }

	// dialer := &kafka.Dialer{
	// 	Timeout:   10 * time.Second,
	// 	TLS:       tlsConfig,
	// 	DualStack: true,
	// }

	// brokers := strings.Split(config.BootstrapServers, ",")
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{config.BootstrapServers},
		Topic:   "reaction-emoji-submission",
		// Dialer:  dialer,
	})
	defer writer.Close()

	server := NewServer(writer)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Service is running")
	})
	http.HandleFunc("/health", server.healthCheckHandler)
	http.HandleFunc("/reaction", server.reactionHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

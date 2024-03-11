package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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
func (s *Server) signHandler(w http.ResponseWriter, r *http.Request) {
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

	err = s.KafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte(req.Message),
		},
	)
	if err != nil {
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
	config, err := LoadKafkaConfig("../../kafka-config.json")
	if err != nil {
		panic(err)
	}

	brokers := strings.Split(config.BootstrapServers, ",")
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   "testtopic1",
	})
	defer writer.Close()

	server := NewServer(writer)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Service is running")
	})

	http.HandleFunc("/health", server.healthCheckHandler)

	http.HandleFunc("/sign", server.signHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type HealthCheckResponse struct {
	Status   string `json:"status"`
	Producer string `json:"producer"`
	Consumer string `json:"consumer"`
}

type SignRequest struct {
	Message string `json:"message"`
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
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

func signHandler(w http.ResponseWriter, r *http.Request) {
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

	signedMessage := req.Message + "-signed"
	fmt.Println("Signed Message:", signedMessage)

	response := map[string]string{"signedMessage": signedMessage}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Service is running")
	})

	http.HandleFunc("/health", healthCheckHandler)

	http.HandleFunc("/sign", signHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

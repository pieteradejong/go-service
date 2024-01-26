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

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthCheckResponse{
		Status:   "Healthy",
		Producer: "Running",
		Consumer: "Running",
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

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Service is running")
	})

	http.HandleFunc("/health", healthCheckHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

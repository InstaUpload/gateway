package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func SendJsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	// NOTE: Is this not working I don't see any message when I call the create api.
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		log.Println("error sending response: ", err)
		return
	}
}

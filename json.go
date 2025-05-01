package main

import (
	"encoding/json"
	"net/http"
)

func SendJsonResponse(w *http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(*w).Encode(data); err != nil {
		http.Error(*w, "Failed to encode response", http.StatusInternalServerError)
	}
}

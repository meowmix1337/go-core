package http_util

import (
	"encoding/json"
	"net/http"
)

// JSONResponse writes a JSON response with the provided status code and data
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	// Set content type header
	w.Header().Set("Content-Type", "application/json")

	// Set status code
	w.WriteHeader(statusCode)

	// Marshal the data to JSON
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If an error occurs during encoding, send an internal server error response
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

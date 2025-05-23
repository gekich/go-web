package response

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, message string, statusCode int) {
	WriteJSON(w, statusCode, map[string]string{
		"error": message,
	})
}

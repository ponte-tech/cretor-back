package response

import (
	"encoding/json"
	"net/http"
)

type ErrorBody struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

type MessageBody struct {
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, message string, status int) {
	JSON(w, ErrorBody{Error: message}, status)
}

func ValidationError(w http.ResponseWriter, err error) {
	JSON(w, ErrorBody{Error: "validation failed", Details: err.Error()}, http.StatusBadRequest)
}

func Message(w http.ResponseWriter, message string, status int) {
	JSON(w, MessageBody{Message: message}, status)
}

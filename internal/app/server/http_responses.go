package server

import (
	"encoding/json"
	"net/http"
	"sports-news-api/internal/app/transport"
	"time"
)

type FailureResponse struct {
	Status   string             `json:"status"`
	Data     interface{}        `json:"data"`
	Metadata transport.Metadata `json:"Metadata"`
}

type ErrorResponse struct {
	Status   string             `json:"status"`
	Message  string             `json:"message"`
	Metadata transport.Metadata `json:"Metadata"`
}

func RespondWithSuccess(data any, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}

func RespondWithFailure(data any, statusCode int, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	response := FailureResponse{
		Status: StatusFail,
		Data:   data,
		Metadata: transport.Metadata{
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		},
	}
	_ = json.NewEncoder(w).Encode(response)
}

func RespondWithError(message string, statusCode int, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	response := ErrorResponse{
		Status:  StatusError,
		Message: message,
		Metadata: transport.Metadata{
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		},
	}
	_ = json.NewEncoder(w).Encode(response)
}

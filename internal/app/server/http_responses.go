package server

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type FailureResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func RespondWithSuccess(data any, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	response := SuccessResponse{
		Status: "success",
		Data:   data,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func RespondWithFailure(data any, statusCode int, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	response := FailureResponse{
		Status: "fail",
		Data:   data,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func RespondWithError(message string, statusCode int, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	response := ErrorResponse{
		Status:  "error",
		Message: message,
	}
	_ = json.NewEncoder(w).Encode(response)
}

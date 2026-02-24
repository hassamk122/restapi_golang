package utils

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func RespondWithSuccess(res http.ResponseWriter, statusCode int, message string, data interface{}) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	json.NewEncoder(res).Encode(SuccessResponse{
		Message: message,
		Data:    data,
	})
}

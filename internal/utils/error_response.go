package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func RespondWithError(res http.ResponseWriter, statusCode int, message string) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	json.NewEncoder(res).Encode(ErrorResponse{
		Message: message,
	})
}

func ResondWithNotFound(res http.ResponseWriter) {
	RespondWithError(res, http.StatusNotFound, "Resource not found")
}

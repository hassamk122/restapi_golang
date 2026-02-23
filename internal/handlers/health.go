package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) HealthHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"message": "server is OK",
		}
		json.NewEncoder(res).Encode(response)
	}
}

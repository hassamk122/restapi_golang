package routes

import (
	"net/http"

	"github.com/hassamk122/restapi_golang/internal/handlers"
)

func SetupHealthRoute(mux *http.ServeMux, handler *handlers.Handler) {
	mux.HandleFunc("/health", handler.HealthHandler())
}

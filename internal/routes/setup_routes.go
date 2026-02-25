package routes

import (
	"net/http"

	"github.com/hassamk122/restapi_golang/internal/handlers"
)

func SetupRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	SetupHealthRoute(mux, handler)
	SetupUserRoutes(mux, handler)
}

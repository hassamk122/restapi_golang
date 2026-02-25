package routes

import (
	"net/http"

	"github.com/hassamk122/restapi_golang/internal/handlers"
)

func SetupUserRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	mux.HandleFunc("POST /user/register", handler.CreateUserHandler())
}

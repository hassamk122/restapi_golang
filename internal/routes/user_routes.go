package routes

import (
	"net/http"

	"github.com/hassamk122/restapi_golang/internal/handlers"
)

func SetupUserRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	userMux := http.NewServeMux()

	mux.Handle("/user/", http.StripPrefix("/user", userMux))

	userMux.HandleFunc("POST /register", handler.CreateUserHandler())
	userMux.HandleFunc("POST /login", handler.LoginUserHandler())
}

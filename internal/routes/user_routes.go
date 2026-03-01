package routes

import (
	"net/http"

	"github.com/hassamk122/restapi_golang/internal/handlers"
	"github.com/hassamk122/restapi_golang/internal/middlewares"
)

func SetupUserRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	userMux := http.NewServeMux()

	mux.Handle("/users/", http.StripPrefix("/users", userMux))

	userMux.Handle("GET /all", middlewares.Apply(
		handler.ListUserHandler(),
	))

	userMux.Handle("POST /register", middlewares.Apply(
		handler.CreateUserHandler(),
	))

	userMux.Handle("POST /login", middlewares.Apply(
		handler.LoginUserHandler(),
	))

	userMux.Handle("POST /session/logout", middlewares.Apply(
		handler.LogoutHandler(),
		middlewares.AuthMiddleware,
	))

	userMux.Handle("GET /current-user/profile", middlewares.Apply(
		handler.UserProfileHandler(),
		middlewares.AuthMiddleware,
	))

	mux.Handle("POST /upload/profile-image", middlewares.Apply(
		handler.UploadProfileImageHandler(),
		middlewares.AuthMiddleware,
	))

}

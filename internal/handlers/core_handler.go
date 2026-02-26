package handlers

import (
	"github.com/hassamk122/restapi_golang/internal/service"
)

type Handler struct {
	UserService *service.UserService
}

func NewHandler(userService *service.UserService) *Handler {
	return &Handler{
		UserService: userService,
	}
}

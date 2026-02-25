package handlers

import (
	"github.com/hassamk122/restapi_golang/internal/repo"
	"github.com/hassamk122/restapi_golang/internal/service"
)

type Handler struct {
	UserRepo    repo.UserRepo
	UserService *service.UserService
}

func NewHandler(userRepo repo.UserRepo, userService *service.UserService) *Handler {
	return &Handler{
		UserRepo:    userRepo,
		UserService: userService,
	}
}

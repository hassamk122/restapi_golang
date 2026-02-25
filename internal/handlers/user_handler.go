package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hassamk122/restapi_golang/internal/dtos"
	"github.com/hassamk122/restapi_golang/internal/store"
	"github.com/hassamk122/restapi_golang/internal/utils"
)

func (h *Handler) CreateUserHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		var userReq dtos.CreateUserRequest
		if err := json.NewDecoder(req.Body).Decode(&userReq); err != nil {
			utils.RespondWithError(res, http.StatusBadGateway, "Invalid request payload")
			return
		}

		hashedPassword, err := utils.HashPassword(userReq.Password)
		if err != nil {
			utils.RespondWithError(res, http.StatusInternalServerError, "Error while hashing password")
			return
		}

		_, err = h.Queries.CreateUser(ctx, store.CreateUserParams{
			Username: userReq.Username,
			Email:    userReq.Email,
			Password: hashedPassword,
		})
		if err != nil {
			utils.RespondWithError(res, http.StatusInternalServerError, "Error creating user")
			return
		}

		utils.RespondWithSuccess(res, http.StatusCreated, "User created successfully", userReq.Username)
	}
}

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hassamk122/restapi_golang/internal/dtos"
	"github.com/hassamk122/restapi_golang/internal/store"
	"github.com/hassamk122/restapi_golang/internal/utils"
	"github.com/hassamk122/restapi_golang/internal/validation"
)

func (h *Handler) CreateUserHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		var userReq dtos.CreateUserRequest
		if err := json.NewDecoder(req.Body).Decode(&userReq); err != nil {
			utils.RespondWithError(res, http.StatusBadGateway, "Invalid request payload")
			return
		}

		if err := validation.Validate(&userReq); err != nil {
			utils.RespondWithError(res, http.StatusBadRequest, err.Error())
			return
		}

		tx, err := h.DB.BeginTx(ctx, nil)
		if err != nil {
			utils.RespondWithError(res, http.StatusBadRequest, "Error starting transaction")
			return
		}
		defer tx.Rollback()

		qtx := store.New(tx)

		_, err = qtx.GetUserByEmail(ctx, userReq.Email)
		if err != nil {
			utils.RespondWithError(res, http.StatusConflict, "Email already taken")
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

		if err := tx.Commit(); err != nil {
			utils.RespondWithError(res, http.StatusInternalServerError, "Failed to commit transaction")
		}

		utils.RespondWithSuccess(res, http.StatusCreated, "User created successfully", userReq.Username)
	}
}

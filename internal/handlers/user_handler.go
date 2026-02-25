package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hassamk122/restapi_golang/internal/dtos"
	"github.com/hassamk122/restapi_golang/internal/utils"
	"github.com/hassamk122/restapi_golang/internal/validation"
)

func (h *Handler) LoginUserHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		var userReq dtos.LoginRequest

		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&userReq); err != nil {
			utils.RespondWithError(res, http.StatusBadGateway, "Invalid request payload")
			return
		}

		if err := validation.Validate(&userReq); err != nil {
			utils.RespondWithError(res, http.StatusBadRequest, err.Error())
			return
		}

		token, err := h.UserService.Login(ctx, userReq.Email, userReq.Password)
		if err != nil {
			utils.RespondWithError(res, http.StatusInternalServerError, "Error generating a token")
			return
		}

		utils.RespondWithSuccess(res, http.StatusOK, "Login sucessful", map[string]string{
			"token": token,
		})
	}
}

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

		err := h.UserService.Register(ctx, userReq.Username, userReq.Email, userReq.Password)
		if errors.Is(err, utils.ErrEmailTaken) {
			utils.RespondWithError(res, http.StatusConflict, "Email already taken")
			return
		}

		if err != nil {
			utils.RespondWithError(res, http.StatusInternalServerError, "Error creating user")
			return
		}

		utils.RespondWithSuccess(res, http.StatusCreated, "User created successfully", nil)
	}
}

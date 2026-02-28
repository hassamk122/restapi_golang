package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/hassamk122/restapi_golang/internal/auth"
	"github.com/hassamk122/restapi_golang/internal/dtos"
	"github.com/hassamk122/restapi_golang/internal/middlewares"
	"github.com/hassamk122/restapi_golang/internal/service"
	"github.com/hassamk122/restapi_golang/internal/utils"
	"github.com/hassamk122/restapi_golang/internal/validation"
)

func (h *Handler) UserProfile() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		claims, ok := req.Context().Value(middlewares.UserClaimsKey).(*auth.Claims)
		if !ok {
			utils.RespondWithError(res, http.StatusBadRequest, "please login to continue")
			return
		}

		userID := claims.UserId

		user, err := h.UserService.GetUserProfile(req.Context(), int32(userID))
		if err != nil {
			utils.RespondWithError(res, http.StatusNotFound, err.Error())
		}

		utils.RespondWithSuccess(res, http.StatusOK, "success", user)
	}
}

func (h *Handler) UploadProfileImageHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		claims, ok := req.Context().Value(middlewares.UserClaimsKey).(*auth.Claims)
		if !ok {
			utils.RespondWithError(res, http.StatusBadRequest, "please login to continue")
			return
		}

		userID := claims.UserId

		file, fileHeader, err := h.UserService.ParseAndRetrieveProfile(req)
		if err != nil {
			if strings.Contains(err.Error(), "parse_multipart_err") {
				utils.RespondWithError(res, http.StatusBadRequest, "Error parsing file")
			} else if strings.Contains(err.Error(), "retrieve_file_err") {
				utils.RespondWithError(res, http.StatusBadRequest, "Error retrieving file")

			}
			return
		}

		uploadedResultURL, err := h.UserService.UploadAndSaveImage(int32(userID), file, fileHeader, req.Context())
		if err != nil {
			if strings.Contains(err.Error(), "cloud_init_failed") {
				utils.RespondWithError(res, http.StatusInternalServerError, "Cloud initialization failed")
			} else if strings.Contains(err.Error(), "upload_failed") {
				utils.RespondWithError(res, http.StatusBadGateway, "Image upload failed")
			}
			return
		}

		utils.RespondWithSuccess(res, http.StatusOK, "User Profile uploaded successfully", uploadedResultURL)
	}
}

func (h *Handler) LoginUserHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		var userReq dtos.LoginRequest

		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&userReq); err != nil {
			utils.RespondWithError(res, http.StatusBadGateway, service.ErrInvalidRequestPayload.Error())
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
			utils.RespondWithError(res, http.StatusBadGateway, service.ErrInvalidRequestPayload.Error())
			return
		}

		if err := validation.Validate(&userReq); err != nil {
			utils.RespondWithError(res, http.StatusBadRequest, err.Error())
			return
		}

		err := h.UserService.Register(ctx, userReq.Username, userReq.Email, userReq.Password)
		if errors.Is(err, service.ErrEmailTaken) {
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

func (h *Handler) LogoutHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		claims, ok := req.Context().Value(middlewares.UserClaimsKey).(*auth.Claims)
		if !ok {
			utils.RespondWithError(res, http.StatusBadRequest, "Please login to continue")
			return
		}

		tokenString := extractTokenFromHeader(req)
		if tokenString == "" {
			utils.RespondWithError(res, http.StatusUnauthorized, "Missing token")
			return
		}

		err := h.UserService.Logout(req.Context(), claims, tokenString)
		if err != nil {
			utils.RespondWithError(res, http.StatusInternalServerError, "Error Logging out")
			return
		}

		utils.RespondWithSuccess(res, http.StatusCreated, "Logged out successfully", true)
	}
}

func extractTokenFromHeader(req *http.Request) string {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

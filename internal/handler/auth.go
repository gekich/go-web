package handler

import (
	"encoding/json"
	"errors"
	dto "github.com/gekich/go-web/internal/dto/auth"
	"github.com/gekich/go-web/internal/service"
	"github.com/gekich/go-web/internal/util"
	message "github.com/gekich/go-web/internal/util/message"
	"github.com/gekich/go-web/internal/util/response"
	"net/http"
)

type AuthHandler struct {
	service   *service.AuthService
	validator util.Validator
}

func NewAuthHandler(s *service.AuthService, v util.Validator) *AuthHandler {
	return &AuthHandler{
		service:   s,
		validator: v,
	}
}

// @Summary Register
// @Description Register a user
// @Accept json
// @Produce json
// @Param User body dto.RegisterUserInput true "Register a new user"
// @Success 200 {string} OK
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @router /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input dto.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, "invalid input", http.StatusBadRequest)
		return
	}

	err := h.validator.Validate(input)
	if err != nil {
		response.WriteError(w, "validation error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.service.Register(r.Context(), input)

	switch {
	case errors.Is(message.ErrEmailInUse, err):
		response.WriteError(w, err.Error(), http.StatusBadRequest)
	case err != nil:
		response.WriteError(w, "registration failed: ", http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
	}
}

// @Summary Login
// @Description Login a user
// @Accept json
// @Produce json
// @Param User body dto.LoginInput true "Login user"
// @Success 200 {string} OK
// @Failure 400 {string} Bad Request
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input dto.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, "invalid input", http.StatusBadRequest)
		return
	}

	err := h.validator.Validate(input)
	if err != nil {
		response.WriteError(w, "validation error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := h.service.Login(r.Context(), input)

	switch {
	case errors.Is(err, message.ErrInvalidCredentials):
		response.WriteError(w, err.Error(), http.StatusUnauthorized)
	case err != nil:
		response.WriteError(w, "login failed", http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))
	}
}

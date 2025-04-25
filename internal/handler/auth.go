package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	dto "github.com/gekich/go-web/internal/dto/auth"
	"github.com/gekich/go-web/internal/service"
	message "github.com/gekich/go-web/internal/util/message"
	"github.com/gekich/go-web/internal/validator"
	"net/http"
)

type AuthHandler struct {
	service   *service.AuthService
	validator validator.Validator
}

func NewAuthHandler(s *service.AuthService, v validator.Validator) *AuthHandler {
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
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	err := h.validator.Validate(input)
	if err != nil {
		http.Error(w, "validation error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.service.Register(r.Context(), input)
	fmt.Println(err)
	fmt.Println(message.ErrEmailInUse)

	switch {
	case errors.Is(message.ErrEmailInUse, err):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case err != nil:
		http.Error(w, "registration failed: ", http.StatusInternalServerError)
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
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	err := h.validator.Validate(input)
	if err != nil {
		http.Error(w, "validation error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := h.service.Login(r.Context(), input)

	switch {
	case errors.Is(err, message.ErrInvalidCredentials):
		http.Error(w, err.Error(), http.StatusUnauthorized)
	case err != nil:
		http.Error(w, "login failed", http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))
	}
}

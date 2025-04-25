package handler

import (
	"encoding/json"
	dto "github.com/gekich/go-web/internal/dto/user"
	service "github.com/gekich/go-web/internal/service/user"
	"github.com/gekich/go-web/internal/validator"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

type UserHandler struct {
	service   *service.UserService
	validator validator.Validator
}

func NewUserHandler(s *service.UserService, v validator.Validator) *UserHandler {
	return &UserHandler{
		service:   s,
		validator: v,
	}
}

// Get an user by its ID
// @Summary Get an User
// @Description Get an user by its id.
// @Accept json
// @Produce json
// @Param id path int true "user ID"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @router /users/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetPublicUser(r.Context(), id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// ListUsers GET /users
// @Summary Shows all users
// @Description Lists all users.
// @Success 200 {array}  dto.UserResponse
// @Failure 500 {string} Internal Server Error
// @router /users [get]
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		http.Error(w, "error fetching users", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, users)
}

// CreateUser creates a new user
// Create
// @Summary Create a User
// @Description Create a user using JSON payload
// @Accept json
// @Produce json
// @Param User body dto.CreateUserInput true "Create a user using the following format"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @router /users/{id} [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	err := h.validator.Validate(input)
	if err != nil {
		http.Error(w, "validation error:"+err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := h.service.CreateUser(r.Context(), input)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

// UpdateUser update user
// Update
// @Summary Update a User
// @Description Update a user using JSON payload
// @Accept json
// @Produce json
// @Param User body dto.UpdateUserInput true "Update user using the following format"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @router /users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var input dto.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateUser(r.Context(), id, input)
	if err != nil {
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	user, err := h.service.GetPublicUser(r.Context(), id)

	writeJSON(w, http.StatusOK, user)
}

// DeleteUser a user by its ID
// @Summary Delete a User
// @Description Delete a user by its id.
// @Accept json
// @Produce json
// @Param id path int true "user ID"
// @Success 200 "Ok"
// @Failure 500 {string} Internal Server Error
// @router /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteUser(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// writeJSON is a helper function
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

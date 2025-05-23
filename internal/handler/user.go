package handler

import (
	"encoding/json"
	dto "github.com/gekich/go-web/internal/dto/user"
	"github.com/gekich/go-web/internal/logger"
	"github.com/gekich/go-web/internal/service"
	"github.com/gekich/go-web/internal/util"
	"github.com/gekich/go-web/internal/util/response"
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
	validator util.Validator
	logger    logger.Logger
}

func NewUserHandler(s *service.UserService, v util.Validator, l logger.Logger) *UserHandler {
	return &UserHandler{
		service:   s,
		validator: v,
		logger:    l,
	}
}

// GetUserByID Get a user by its ID
// @Summary Get a User
// @Description Get a user by its id.
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
		response.WriteError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetPublicUser(r.Context(), id)
	if err != nil {
		response.WriteError(w, "user not found", http.StatusNotFound)
		return
	}

	response.WriteJSON(w, http.StatusOK, user)
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
		response.WriteError(w, "error fetching users", http.StatusInternalServerError)
		return
	}

	response.WriteJSON(w, http.StatusOK, users)
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
		response.WriteError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var input dto.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, "invalid input", http.StatusBadRequest)
		return
	}

	err = h.validator.Validate(input)
	if err != nil {
		response.WriteError(w, "validation error:"+err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.service.UpdateUser(r.Context(), id, input)
	if err != nil {
		response.WriteError(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	user, err := h.service.GetPublicUser(r.Context(), id)

	response.WriteJSON(w, http.StatusOK, user)
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
		response.WriteError(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteUser(r.Context(), id)
	if err != nil {
		response.WriteError(w, "failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

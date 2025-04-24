package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	db "github.com/gekich/go-web/internal/db/user"
	"github.com/go-chi/chi/v5"
)

type UserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

type UserHandler struct {
	Queries *db.Queries
}

func NewUserHandler(queries *db.Queries) *UserHandler {
	return &UserHandler{Queries: queries}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, err := h.Queries.GetUser(r.Context(), int64(id))

	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	res := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

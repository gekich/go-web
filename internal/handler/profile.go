package handler

import (
	"github.com/gekich/go-web/internal/middleware"
	"github.com/gekich/go-web/internal/service"
	"github.com/gekich/go-web/internal/util/response"
	"net/http"
	"strconv"
)

type ProfileHandler struct {
	service *service.UserService
}

func NewProfileHandler(s *service.UserService) *ProfileHandler {
	return &ProfileHandler{
		service: s,
	}
}

// @Summary Profile
// @Description Get user profile
// @Accept json
// @Produce json
// @Param User body dto.UserResponse true "Get user information"
// @Success 200 {string} OK
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @router /profile [get]
func (h *ProfileHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.UserIDFromContext(r.Context())

	if !ok {
		response.WriteError(w, "user unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		response.WriteError(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetPublicUser(r.Context(), id)
	if err != nil {
		response.WriteError(w, "user not found", http.StatusNotFound)
		return
	}

	response.WriteJSON(w, http.StatusOK, user)
}

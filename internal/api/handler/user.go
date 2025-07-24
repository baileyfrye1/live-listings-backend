package handler

import (
	"encoding/json"
	"net/http"

	"server/internal/api/dto"
	"server/internal/server/middleware"
	"server/internal/service"
	"server/util"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetAllAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := h.userService.GetUsersByRole(r.Context(), "agent")
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Could not find agents")
		return
	}

	util.WriteJSON(w, http.StatusOK, agents)
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserContextKey).(int)

	user, err := h.userService.GetUserById(r.Context(), userId)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Could not find user")
		return
	}

	util.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUserById(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserContextKey).(int)
	var req dto.RequestUpdateUser

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid fields")
		return
	}

	user, err := h.userService.UpdateUserById(r.Context(), &req, userId)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Could not find user")
		return
	}

	util.WriteJSON(w, http.StatusOK, user)
}

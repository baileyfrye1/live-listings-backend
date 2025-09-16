package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"server/internal/api/dto"
	"server/internal/domain"
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

func (h *UserHandler) GetAgentById(w http.ResponseWriter, r *http.Request) {
	agentId, err := strconv.Atoi(chi.URLParam(r, "agentId"))
	if err != nil {
		util.RespondWithError(
			w,
			http.StatusBadRequest,
			"Incorrect ID provided. Please provide a valid ID",
		)
		return
	}

	agent, err := h.userService.GetAgentById(r.Context(), agentId)
	if err != nil {
		util.RespondWithError(
			w,
			http.StatusInternalServerError,
			"Agent could not be found with that ID",
		)
		return
	}

	util.WriteJSON(w, http.StatusOK, agent)
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(middleware.UserContextKey).(*domain.ContextSessionData)

	user, err := h.userService.GetUserById(r.Context(), userCtx.UserID)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Could not find user")
		return
	}

	util.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUserById(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(middleware.UserContextKey).(*domain.ContextSessionData)
	var req dto.UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid fields")
		return
	}

	user, err := h.userService.UpdateUserById(r.Context(), &req, userCtx.UserID)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Error updating user")
		return
	}

	util.WriteJSON(w, http.StatusOK, user)
}

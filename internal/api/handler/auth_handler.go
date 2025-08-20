package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"server/internal/api/dto"
	"server/internal/domain"
	"server/internal/server/middleware"
	"server/internal/service"
	"server/util"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	user, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		slog.Error("Create User - Service Error", slog.String("error", err.Error()))
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    user.SessionID,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	util.WriteJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid email or password")
		return
	}

	user, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    user.SessionID,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "Successfully logged in"})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(middleware.UserContextKey).(*domain.ContextSessionData)
	err := h.authService.Logout(r.Context(), userCtx.SessionID)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Expires:  time.Now().Add(-1),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "Successfully logged out"})
}

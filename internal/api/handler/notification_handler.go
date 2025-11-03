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

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

func (h *NotificationHandler) GetAllNotificationsByUserId(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(middleware.UserContextKey).(*domain.ContextSessionData)

	notifications, err := h.notificationService.GetAllNotificationsByUserId(
		r.Context(),
		userCtx.UserID,
	)
	if err != nil {
		util.RespondWithError(
			w,
			http.StatusInternalServerError,
			"Error fetching user notifications",
		)
		return
	}

	util.WriteJSON(w, http.StatusOK, notifications)
}

func (h *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(middleware.UserContextKey).(*domain.ContextSessionData)

	var req dto.CreateNotificationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	newNotification := &domain.Notification{
		UserID:    userCtx.UserID,
		ListingID: req.ListingID,
		Type:      req.Type,
		Message:   req.Message,
		IsRead:    false,
	}

	notification, err := h.notificationService.CreateNotification(r.Context(), newNotification)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Error creating notification")
		return
	}

	util.WriteJSON(w, http.StatusOK, notification)
}

func (h *NotificationHandler) ToggleNotificationReadStatus(w http.ResponseWriter, r *http.Request) {
	notificationId, err := strconv.Atoi(chi.URLParam(r, "notificationId"))
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid id")
		return
	}

	notification, err := h.notificationService.ToggleNotificationReadStatus(
		r.Context(),
		notificationId,
	)
	if err != nil {
		util.RespondWithError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	util.WriteJSON(w, http.StatusOK, notification)
}

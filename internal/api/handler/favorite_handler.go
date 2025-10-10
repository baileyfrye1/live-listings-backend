package handler

import (
	"encoding/json"
	"net/http"

	"server/internal/api/dto"
	"server/internal/domain"
	"server/internal/server/middleware"
	"server/internal/service"
	"server/util"
)

type FavoriteHandler struct {
	favoriteService *service.FavoriteService
}

func NewFavoriteHandler(favoriteService *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{favoriteService: favoriteService}
}

func (h *FavoriteHandler) CreateFavorite(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(middleware.UserContextKey).(*domain.ContextSessionData)
	var req dto.FavoriteListingDto

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	newFavorite := &domain.Favorite{UserID: userCtx.UserID, ListingID: req.ListingID}

	favorite, err := h.favoriteService.CreateFavorite(r.Context(), newFavorite)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Error creating favorite")
		return
	}

	util.WriteJSON(w, http.StatusOK, favorite)
}

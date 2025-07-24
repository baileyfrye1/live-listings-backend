package handler

import (
	"net/http"

	"server/internal/service"
	"server/util"
)

type ListingHandler struct {
	listingService *service.ListingService
}

func NewListingHandler(listingService *service.ListingService) *ListingHandler {
	return &ListingHandler{listingService: listingService}
}

func (h *ListingHandler) GetAllListings(w http.ResponseWriter, r *http.Request) {
	listings, err := h.listingService.GetAllListings(r.Context())
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	util.WriteJSON(w, http.StatusOK, listings)
}

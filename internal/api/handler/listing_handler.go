package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"server/internal/api/dto"
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

func (h *ListingHandler) GetListingById(w http.ResponseWriter, r *http.Request) {
	listingId, err := strconv.Atoi(chi.URLParam(r, "listingId"))
	if err != nil {
		util.RespondWithError(
			w,
			http.StatusBadRequest,
			"Listing id is in incorrect format",
		)
		return
	}

	listing, err := h.listingService.GetListingById(r.Context(), listingId)
	if err != nil {
		util.RespondWithError(
			w,
			http.StatusBadRequest,
			"Listing could not be found",
		)
		return
	}

	util.WriteJSON(w, http.StatusOK, listing)
}

func (h *ListingHandler) CreateListing(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestCreateListing

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Please enter all fields")
		return
	}

	listing, err := h.listingService.CreateListing(r.Context(), &req)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, listing)
}

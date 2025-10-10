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

type ListingHandler struct {
	listingService *service.ListingService
	userService    *service.UserService
}

func NewListingHandler(
	listingService *service.ListingService,
	userService *service.UserService,
) *ListingHandler {
	return &ListingHandler{listingService: listingService, userService: userService}
}

func (h *ListingHandler) GetAllListings(w http.ResponseWriter, r *http.Request) {
	listings, err := h.listingService.GetAllListings(r.Context())
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	util.WriteJSON(w, http.StatusOK, listings)
}

func (h *ListingHandler) GetMyListings(w http.ResponseWriter, r *http.Request) {
	currentAgentCtx := r.Context().Value(middleware.UserContextKey).(*domain.ContextSessionData)

	listings, err := h.listingService.GetListingsByAgentId(r.Context(), currentAgentCtx.UserID)
	if err != nil {
		util.RespondWithError(
			w,
			http.StatusInternalServerError,
			"Could not fetch current agent listings",
		)
		return
	}

	if len(listings) == 0 {
		util.WriteJSON(w, http.StatusOK, []domain.Listing{})
		return
	}

	util.WriteJSON(w, http.StatusOK, listings)
}

func (h *ListingHandler) GetAgentListings(w http.ResponseWriter, r *http.Request) {
	agentId, err := strconv.Atoi(chi.URLParam(r, "agentId"))
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Agent id is in incorrect format")
		return
	}

	_, err = h.userService.GetUserById(r.Context(), agentId)
	if err != nil {
		util.RespondWithError(w, http.StatusNotFound, "Agent does not exist")
		return
	}

	listings, err := h.listingService.GetListingsByAgentId(r.Context(), agentId)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Could not fetch listings")
		return
	}

	if len(listings) == 0 {
		util.WriteJSON(w, http.StatusOK, []domain.Listing{})
		return
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
	agentCtx := r.Context().Value(middleware.UserContextKey).(*domain.ContextSessionData)
	var req dto.CreateListingRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Please enter all fields")
		return
	}

	if (agentCtx.Role == "agent") && req.AgentID != nil {
		util.RespondWithError(
			w,
			http.StatusBadRequest,
			"Cannot change agent on listing. Please contact admin to change agent",
		)
		return
	}

	if req.AgentID == nil {
		req.AgentID = &agentCtx.UserID
	}

	listing, err := h.listingService.CreateListing(r.Context(), &req)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, listing)
}

func (h *ListingHandler) UpdateMyListing(w http.ResponseWriter, r *http.Request) {
	currentUserCtx := r.Context().Value(middleware.UserContextKey).(*domain.ContextSessionData)
	listingId, err := strconv.Atoi(chi.URLParam(r, "listingId"))
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Incorrect ID format")
	}

	var req dto.UpdateListingRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Please enter a field to update")
		return
	}

	listing, err := h.listingService.UpdateListingById(
		r.Context(),
		&req,
		currentUserCtx,
		listingId,
	)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, listing)
}

func (h *ListingHandler) DeleteMyListing(w http.ResponseWriter, r *http.Request) {
	currentUserCtx := r.Context().Value(middleware.UserContextKey).(*domain.ContextSessionData)
	listingId, err := strconv.Atoi(chi.URLParam(r, "listingId"))
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Incorrect ID format")
		return
	}

	err = h.listingService.DeleteListingById(r.Context(), currentUserCtx, listingId)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "Successfully deleted"})
}

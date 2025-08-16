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
	currentAgentId := r.Context().Value(middleware.UserContextKey).(int)

	listings, err := h.listingService.GetListingsByAgentId(r.Context(), currentAgentId)
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
	userRole := r.Context().Value(middleware.RoleContextKey).(string)
	userId := r.Context().Value(middleware.UserContextKey).(int)
	var req dto.RequestCreateListing

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Please enter all fields")
		return
	}

	if (userRole == "agent") && req.AgentID != nil {
		util.RespondWithError(
			w,
			http.StatusBadRequest,
			"Cannot change agent on listing. Please contact admin to change agent",
		)
		return
	}

	if req.AgentID == nil {
		req.AgentID = &userId
	}

	listing, err := h.listingService.CreateListing(r.Context(), &req)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, listing)
}

package service

import (
	"context"
	"errors"

	"server/internal/api/dto"
	"server/internal/domain"
	"server/internal/repo"
)

type ListingService struct {
	listingRepo *repo.ListingRepository
}

func NewListingService(listingRepo *repo.ListingRepository) *ListingService {
	return &ListingService{listingRepo: listingRepo}
}

func (s *ListingService) GetAllListings(ctx context.Context) ([]*domain.Listing, error) {
	return s.listingRepo.GetAllListings(ctx)
}

func (s *ListingService) GetListingsByAgentId(
	ctx context.Context,
	agentId int,
) ([]*domain.Listing, error) {
	return s.listingRepo.GetListingsByAgentId(ctx, agentId)
}

func (s *ListingService) GetListingById(ctx context.Context, id int) (*domain.Listing, error) {
	return s.listingRepo.GetListingById(ctx, id)
}

func (s *ListingService) CreateListing(
	ctx context.Context,
	listing *dto.CreateListingRequest,
) (*domain.Listing, error) {
	return s.listingRepo.CreateListing(ctx, listing)
}

func (s *ListingService) UpdateListingById(
	ctx context.Context,
	listing *dto.UpdateListingRequest,
	currentUserCtx *domain.ContextSessionData,
	listingId int,
) (*domain.Listing, error) {
	if (currentUserCtx.Role == "agent") && listing.AgentID != nil {
		return nil, errors.New(
			"Cannot update agent on listing. Please contact admin to change agent",
		)
	}

	return s.listingRepo.UpdateListingById(ctx, listing, currentUserCtx, listingId)
}

func (s *ListingService) DeleteListingById(
	ctx context.Context,
	currentUserCtx *domain.ContextSessionData,
	listingId int,
) error {
	return s.listingRepo.DeleteListingById(ctx, currentUserCtx, listingId)
}

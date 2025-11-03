package service

import (
	"context"
	"errors"

	"server/internal/api/dto"
	"server/internal/domain"
	listingRepo "server/internal/repo"
)

type ListingService struct {
	listingRepo listingRepo.IListingRepo
}

func NewListingService(listingRepo listingRepo.IListingRepo) *ListingService {
	return &ListingService{listingRepo: listingRepo}
}

func (s *ListingService) GetAllListings(ctx context.Context) ([]*domain.Listing, error) {
	listings, err := s.listingRepo.GetAllListings(ctx)
	if err != nil {
		return nil, err
	}

	if listings == nil {
		listings = []*domain.Listing{}
	}

	return listings, nil
}

func (s *ListingService) GetListingsByAgentId(
	ctx context.Context,
	agentId int,
) ([]*domain.Listing, error) {
	listings, err := s.listingRepo.GetListingsByAgentId(ctx, agentId)
	if err != nil {
		return nil, err
	}

	if listings == nil {
		listings = []*domain.Listing{}
	}

	return listings, nil
}

func (s *ListingService) GetListingById(ctx context.Context, id int) (*domain.Listing, error) {
	return s.listingRepo.GetListingById(ctx, id)
}

func (s *ListingService) CreateListing(
	ctx context.Context,
	listing *domain.Listing,
) (*domain.Listing, error) {
	return s.listingRepo.CreateListing(ctx, listing)
}

func (s *ListingService) UpdateListingById(
	ctx context.Context,
	listingReq *dto.UpdateListingRequest,
	currentUserCtx *domain.ContextSessionData,
	listingId int,
) (*domain.Listing, error) {
	if (currentUserCtx.Role == "agent") && listingReq.AgentID != nil {
		return nil, errors.New(
			"Cannot update agent on listing. Please contact admin to change agent",
		)
	}

	return s.listingRepo.UpdateListingById(ctx, listingReq, currentUserCtx, listingId)
}

func (s *ListingService) DeleteListingById(
	ctx context.Context,
	currentUserCtx *domain.ContextSessionData,
	listingId int,
) error {
	return s.listingRepo.DeleteListingById(ctx, currentUserCtx, listingId)
}

func (s *ListingService) TrackViewsByListingId(
	ctx context.Context,
	listingId int,
) error {
	return s.listingRepo.TrackViewsByListingId(ctx, listingId)
}

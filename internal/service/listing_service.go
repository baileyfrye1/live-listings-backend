package service

import (
	"context"

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
	currentAgentId int,
	listingId int,
) (*domain.Listing, error) {
	return s.listingRepo.UpdateListingById(ctx, listing, currentAgentId, listingId)
}

func (s *ListingService) DeleteListingById(
	ctx context.Context,
	currentAgentId int,
	listingId int,
) error {
	return s.listingRepo.DeleteListingById(ctx, currentAgentId, listingId)
}

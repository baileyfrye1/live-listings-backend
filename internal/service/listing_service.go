package service

import (
	"context"

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

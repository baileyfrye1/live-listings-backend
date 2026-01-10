package service

import (
	"context"
	"errors"
	"mime/multipart"

	"server/internal/api/dto"
	"server/internal/cloudinary"
	"server/internal/domain"
	listingRepo "server/internal/repo"
)

type ListingService struct {
	listingRepo listingRepo.IListingRepo
	cld         *cloudinary.CloudinaryClient
}

func NewListingService(
	listingRepo listingRepo.IListingRepo,
	cld *cloudinary.CloudinaryClient,
) *ListingService {
	return &ListingService{listingRepo: listingRepo, cld: cld}
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
	multiPartForm *multipart.Form,
	listing *domain.Listing,
) (*domain.Listing, error) {
	files := multiPartForm.File["images"]

	// TODO: Convert this to go routine
	for i, fh := range files {
		file, err := fh.Open()
		if err != nil {
			return nil, err
		}

		result, err := s.cld.Upload(ctx, file, listing.ID)
		file.Close()

		if err != nil {
			return nil, err
		}

		listing.Images = append(listing.Images, domain.ListingImage{
			PublicID:  result.PublicID,
			ListingID: listing.ID,
			URL:       result.URL,
			SortOrder: i,
			IsPrimary: i == 0,
		})
	}

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

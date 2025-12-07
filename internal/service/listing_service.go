package service

import (
	"context"
	"errors"
	"mime/multipart"

	"golang.org/x/sync/errgroup"

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

	type uploadResult struct {
		index int
		image domain.ListingImage
		err   error
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(5)

	results := make(chan uploadResult, len(files))

	for i, fh := range files {
		g.Go(func() error {
			file, err := fh.Open()
			if err != nil {
				return err
			}

			result, err := s.cld.Upload(gCtx, file, listing.ID)
			file.Close()

			if err != nil {
				results <- uploadResult{index: i, image: domain.ListingImage{}, err: err}
				return err
			}

			image := domain.ListingImage{
				PublicID:  result.PublicID,
				ListingID: listing.ID,
				URL:       result.URL,
				SortOrder: i,
				IsPrimary: i == 0,
			}

			results <- uploadResult{index: i, image: image, err: nil}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		close(results)
		return nil, err
	}
	close(results)

	sortedResults := make([]uploadResult, len(files))

	for res := range results {
		sortedResults[res.index] = res
	}

	for _, res := range sortedResults {
		if res.err != nil {
			return nil, res.err
		}

		listing.Images = append(listing.Images, res.image)
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

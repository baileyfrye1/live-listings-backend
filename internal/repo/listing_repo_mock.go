package repo

import (
	"context"

	"server/internal/api/dto"
	"server/internal/domain"
)

type ListingRepoMock struct {
	GetAllListingsFunc        func(ctx context.Context) ([]*domain.Listing, error)
	GetListingByIdFunc        func(ctx context.Context, id int) (*domain.Listing, error)
	GetListingsByAgentIdFunc  func(ctx context.Context, agentId int) ([]*domain.Listing, error)
	CreateListingFunc         func(ctx context.Context, listing *domain.Listing) (*domain.Listing, error)
	UpdateListingByIdFunc     func(ctx context.Context, listing *dto.UpdateListingRequest, currentUserCtx *domain.ContextSessionData, listingId int) (*domain.Listing, error)
	DeleteListingByIdFunc     func(ctx context.Context, userCtx *domain.ContextSessionData, listingId int) error
	GetAgentIdByListingIdFunc func(ctx context.Context, listingId int) (int, error)
	TrackViewsByListingIdFunc func(ctx context.Context, listingId int) error
}

func (l *ListingRepoMock) GetAllListings(ctx context.Context) ([]*domain.Listing, error) {
	return l.GetAllListingsFunc(ctx)
}

func (l *ListingRepoMock) GetListingById(ctx context.Context, id int) (*domain.Listing, error) {
	return l.GetListingByIdFunc(ctx, id)
}

func (l *ListingRepoMock) GetListingsByAgentId(
	ctx context.Context,
	agentId int,
) ([]*domain.Listing, error) {
	return l.GetListingsByAgentIdFunc(ctx, agentId)
}

func (l *ListingRepoMock) CreateListing(
	ctx context.Context,
	listing *domain.Listing,
) (*domain.Listing, error) {
	return l.CreateListingFunc(ctx, listing)
}

func (l *ListingRepoMock) UpdateListingById(
	ctx context.Context,
	listingReq *dto.UpdateListingRequest,
	userCtx *domain.ContextSessionData,
	id int,
) (*domain.Listing, error) {
	return l.UpdateListingByIdFunc(ctx, listingReq, userCtx, id)
}

func (l *ListingRepoMock) DeleteListingById(
	ctx context.Context,
	currentUserCtx *domain.ContextSessionData,
	listingId int,
) error {
	return l.DeleteListingByIdFunc(ctx, currentUserCtx, listingId)
}

func (l *ListingRepoMock) GetAgentIdByListingId(ctx context.Context, listingId int) (int, error) {
	return l.GetAgentIdByListingIdFunc(ctx, listingId)
}

func (l *ListingRepoMock) TrackViewsByListingId(ctx context.Context, listingId int) error {
	return l.TrackViewsByListingIdFunc(ctx, listingId)
}

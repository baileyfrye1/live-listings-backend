package repo

import (
	"context"

	"server/internal/domain"
)

type FavoriteRepoMock struct {
	GetUserFavoritesFunc          func(ctx context.Context, userCtx *domain.ContextSessionData) ([]*domain.Favorite, error)
	CreateFavoriteFunc            func(ctx context.Context, favorite *domain.Favorite) (*domain.Favorite, error)
	DeleteFavoriteByListingIdFunc func(ctx context.Context, listingId int, userCtx *domain.ContextSessionData) error
	GetAllUserIdsByListingIdFunc  func(ctx context.Context, listingId int) (map[int]bool, error)
}

func (f *FavoriteRepoMock) GetUserFavorites(
	ctx context.Context,
	userCtx *domain.ContextSessionData,
) ([]*domain.Favorite, error) {
	return f.GetUserFavoritesFunc(ctx, userCtx)
}

func (f *FavoriteRepoMock) CreateFavorite(
	ctx context.Context,
	favorite *domain.Favorite,
) (*domain.Favorite, error) {
	return f.CreateFavoriteFunc(ctx, favorite)
}

func (f *FavoriteRepoMock) DeleteFavoriteByListingId(
	ctx context.Context,
	listingId int,
	userCtx *domain.ContextSessionData,
) error {
	return f.DeleteFavoriteByListingIdFunc(ctx, listingId, userCtx)
}

func (f *FavoriteRepoMock) GetAllUserIdsByListingId(
	ctx context.Context,
	listingId int,
) (map[int]bool, error) {
	return f.GetAllUserIdsByListingIdFunc(ctx, listingId)
}

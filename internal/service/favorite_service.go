package service

import (
	"context"

	"server/internal/domain"
	"server/internal/repo"
)

type FavoriteService struct {
	favoriteRepo repo.IFavoriteRepo
}

func NewFavoriteService(favoriteRepo repo.IFavoriteRepo) *FavoriteService {
	return &FavoriteService{favoriteRepo: favoriteRepo}
}

func (s *FavoriteService) GetUserFavorites(
	ctx context.Context,
	userCtx *domain.ContextSessionData,
) ([]*domain.Favorite, error) {
	favorites, err := s.favoriteRepo.GetUserFavorites(ctx, userCtx)
	if err != nil {
		return nil, err
	}

	if favorites == nil {
		favorites = []*domain.Favorite{}
	}

	return favorites, nil
}

func (s *FavoriteService) CreateFavorite(
	ctx context.Context,
	favorite *domain.Favorite,
) (*domain.Favorite, error) {
	return s.favoriteRepo.CreateFavorite(ctx, favorite)
}

func (s *FavoriteService) DeleteFavoriteByListingId(
	ctx context.Context,
	listingId int,
	userCtx *domain.ContextSessionData,
) error {
	return s.favoriteRepo.DeleteFavoriteByListingId(ctx, listingId, userCtx)
}

func (s *FavoriteService) GetUserFavoritesMap(
	ctx context.Context,
	userCtx *domain.ContextSessionData,
) (map[int]bool, error) {
	favorites, err := s.GetUserFavorites(ctx, userCtx)
	if err != nil {
		return nil, err
	}

	favMap := map[int]bool{}
	for _, fav := range favorites {
		favMap[fav.ListingID] = true
	}

	return favMap, nil
}

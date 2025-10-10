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

func (s *FavoriteService) CreateFavorite(
	ctx context.Context,
	favorite *domain.Favorite,
) (*domain.Favorite, error) {
	return s.favoriteRepo.CreateFavorite(ctx, favorite)
}

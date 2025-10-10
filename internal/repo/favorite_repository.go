package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"server/internal/domain"
)

type IFavoriteRepo interface {
	CreateFavorite(ctx context.Context, favorite *domain.Favorite) (*domain.Favorite, error)
}

type FavoriteRepo struct {
	db *sql.DB
}

func NewFavoriteRepo(db *sql.DB) *FavoriteRepo {
	return &FavoriteRepo{db: db}
}

func (r *FavoriteRepo) CreateFavorite(
	ctx context.Context,
	favorite *domain.Favorite,
) (*domain.Favorite, error) {
	query := `
		INSERT into favorites (user_id, listing_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, listing_id) DO NOTHING
		RETURNING id, created_at, updated_at
	`

	newFavorite := *favorite
	err := r.db.QueryRowContext(ctx, query, favorite.UserID, favorite.ListingID).
		Scan(&newFavorite.ID, &newFavorite.CreatedAt, &newFavorite.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("Favorite already exists")
	}

	if err != nil {
		return nil, fmt.Errorf("Create favorite: %w", err)
	}

	return &newFavorite, nil
}

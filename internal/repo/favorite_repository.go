package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"server/internal/domain"
)

type IFavoriteRepo interface {
	GetUserFavorites(
		ctx context.Context,
		userCtx *domain.ContextSessionData,
	) ([]*domain.Favorite, error)
	CreateFavorite(ctx context.Context, favorite *domain.Favorite) (*domain.Favorite, error)
	DeleteFavoriteByListingId(
		ctx context.Context,
		listingId int,
		userCtx *domain.ContextSessionData,
	) error
}

type FavoriteRepo struct {
	db *sql.DB
}

func NewFavoriteRepo(db *sql.DB) *FavoriteRepo {
	return &FavoriteRepo{db: db}
}

func (r *FavoriteRepo) GetUserFavorites(
	ctx context.Context,
	userCtx *domain.ContextSessionData,
) ([]*domain.Favorite, error) {
	query := `
		SELECT * FROM favorites
		WHERE user_id = $1
	`

	var favorites []*domain.Favorite

	rows, err := r.db.QueryContext(ctx, query, userCtx.UserID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		favorite := new(domain.Favorite)

		err := rows.Scan(
			&favorite.ID,
			&favorite.UserID,
			&favorite.ListingID,
			&favorite.CreatedAt,
			&favorite.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		favorites = append(favorites, favorite)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return favorites, nil
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

func (r *FavoriteRepo) DeleteFavoriteByListingId(
	ctx context.Context,
	listingId int,
	userCtx *domain.ContextSessionData,
) error {
	query := `
		DELETE FROM favorites
		WHERE listing_id = $1 AND
		(
			user_id = $2 OR
			$3 = 'admin'
		)
	`

	result, err := r.db.ExecContext(ctx, query, listingId, userCtx.UserID, userCtx.Role)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows < 1 {
		return errors.New("Cannot delete other favorites")
	}

	return nil
}

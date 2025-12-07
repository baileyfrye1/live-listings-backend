package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

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
	GetAllUserIdsByListingId(ctx context.Context, listingId int) (map[int]bool, error)
}

type FavoriteRepo struct {
	db *pgxpool.Pool
}

func NewFavoriteRepo(db *pgxpool.Pool) *FavoriteRepo {
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

	rows, err := r.db.Query(ctx, query, userCtx.UserID)
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
	err := r.db.QueryRow(ctx, query, favorite.UserID, favorite.ListingID).
		Scan(&newFavorite.ID, &newFavorite.CreatedAt, &newFavorite.UpdatedAt)

	if err == pgx.ErrNoRows {
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

	result, err := r.db.Exec(ctx, query, listingId, userCtx.UserID, userCtx.Role)
	if err != nil {
		return err
	}

	rows := result.RowsAffected()

	if rows < 1 {
		return errors.New("Cannot delete other favorites")
	}

	return nil
}

func (r *FavoriteRepo) GetAllUserIdsByListingId(
	ctx context.Context,
	listingId int,
) (map[int]bool, error) {
	query := `
		SELECT user_id FROM favorites
		WHERE listing_id = $1
	`

	userIds := make(map[int]bool)

	rows, err := r.db.Query(ctx, query, listingId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var userId int

		if err := rows.Scan(&userId); err != nil {
			return nil, err
		}

		userIds[userId] = true
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userIds, nil
}

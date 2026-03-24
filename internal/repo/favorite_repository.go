package repo

import (
	"context"
	"encoding/json"
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
		SELECT
			f.id,
			f.user_id,
			f.listing_id,
			f.created_at,
			f.updated_at,
			l.id,
			l.address,
			l.price,
			l.beds,
			l.baths,
			l.sq_ft,
			l.description,
			l.agent_id,
			l.created_at,
			l.updated_at,
			l.views,
			u.id AS agent_id,
			u.first_name,
			u.last_name,
			u.email,
			COALESCE(
				json_agg(
					json_build_object(
						'id', li.id,
						'public_id', li.public_id,
						'listing_id', li.listing_id,
						'url', li.url,
						'sort_order', li.sort_order,
						'is_primary', li.is_primary,
						'created_at', li.created_at,
						'updated_at', li.updated_at
					) ORDER BY li.sort_order ASC
				) FILTER (WHERE li.id IS NOT NULL),
				'[]'
			) AS images
		FROM favorites f
		INNER JOIN listings l ON l.id = f.listing_id
		INNER JOIN users u ON l.agent_id = u.id
		LEFT JOIN listing_images li ON l.id = li.listing_id
		WHERE f.user_id = $1
		GROUP BY f.id, l.id, u.id, u.first_name, u.last_name, u.email
	`

	var favorites []*domain.Favorite

	rows, err := r.db.Query(ctx, query, userCtx.UserID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		favorite := new(domain.Favorite)
		listing := new(domain.Listing)
		listing.Agent = new(domain.Agent)
		var imagesJson []byte

		err := rows.Scan(
			&favorite.ID,
			&favorite.UserID,
			&favorite.ListingID,
			&favorite.CreatedAt,
			&favorite.UpdatedAt,
			&listing.ID,
			&listing.Address,
			&listing.Price,
			&listing.Beds,
			&listing.Baths,
			&listing.SqFt,
			&listing.Description,
			&listing.AgentID,
			&listing.CreatedAt,
			&listing.UpdatedAt,
			&listing.Views,
			&listing.Agent.ID,
			&listing.Agent.FirstName,
			&listing.Agent.LastName,
			&listing.Agent.Email,
			&imagesJson,
		)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(imagesJson, &listing.Images); err != nil {
			return nil, err
		}

		favorite.Listing = listing

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

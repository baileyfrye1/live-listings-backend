package repo

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"server/internal/api/dto"
	"server/internal/domain"
)

type IListingRepo interface {
	GetAllListings(ctx context.Context) ([]*domain.Listing, error)
	GetListingById(ctx context.Context, id int) (*domain.Listing, error)
	GetListingsByAgentId(ctx context.Context, agentId int) ([]*domain.Listing, error)
	CreateListing(ctx context.Context, listing *domain.Listing) (*domain.Listing, error)
	UpdateListingById(
		ctx context.Context,
		listing *dto.UpdateListingRequest,
		currentUserCtx *domain.ContextSessionData,
		listingId int,
	) (*domain.Listing, error)
	DeleteListingById(
		ctx context.Context,
		currentUserCtx *domain.ContextSessionData,
		listingId int,
	) error
	GetAgentIdByListingId(ctx context.Context, listingId int) (int, error)
	TrackViewsByListingId(ctx context.Context, listingId int) error
}

type ListingRepository struct {
	db *pgxpool.Pool
}

func NewListingRepository(db *pgxpool.Pool) *ListingRepository {
	return &ListingRepository{db: db}
}

func (r *ListingRepository) GetAllListings(ctx context.Context) ([]*domain.Listing, error) {
	query := `
		SELECT 
		l.*,
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
		FROM listings l
		INNER JOIN users u ON l.agent_id = u.id
		LEFT JOIN listing_images li ON l.id = li.listing_id
		GROUP BY l.id, u.id, u.first_name, u.last_name, u.email;
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var listings []*domain.Listing
	var imagesJson []byte
	for rows.Next() {
		listing := new(domain.Listing)
		listing.Agent = new(domain.Agent)

		err := rows.Scan(
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

		listings = append(listings, listing)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return listings, nil
}

func (r *ListingRepository) GetListingById(ctx context.Context, id int) (*domain.Listing, error) {
	query := `
		SELECT l.*,
			users.id AS agent_id,
			users.first_name,
			users.last_name,
			users.email,
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
		FROM listings l
		INNER JOIN users
			ON l.agent_id = users.id
		LEFT JOIN listing_images li 
			ON l.id = li.listing_id	
		WHERE l.id = $1
		GROUP BY l.id, users.id, users.first_name, users.last_name, users.email
	`

	var listing domain.Listing
	listing.Agent = new(domain.Agent)
	var imagesJson []byte
	err := r.db.QueryRow(ctx, query, id).Scan(
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

	return &listing, nil
}

func (r *ListingRepository) GetListingsByAgentId(
	ctx context.Context,
	agentId int,
) ([]*domain.Listing, error) {
	query := `
		SELECT l.*,
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
		FROM listings l
		INNER JOIN users u
			on l.agent_id = u.id
		LEFT JOIN listing_images li 
			ON l.id = li.listing_id
		WHERE agent_id = $1
		GROUP BY l.id, u.id, u.first_name, u.last_name, u.email;
	`
	rows, err := r.db.Query(ctx, query, agentId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var listings []*domain.Listing
	var imagesJson []byte
	for rows.Next() {
		listing := new(domain.Listing)
		listing.Agent = new(domain.Agent)

		err := rows.Scan(
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

		listings = append(listings, listing)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return listings, nil
}

func (r *ListingRepository) CreateListing(
	ctx context.Context,
	listing *domain.Listing,
) (*domain.Listing, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	batch := &pgx.Batch{}

	// Create listing SQL query
	listingQuery := `
		WITH new_listing AS (
			INSERT INTO listings (address, price, beds, baths, sq_ft, agent_id)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, created_at, updated_at, agent_id
		)
		SELECT nl.*, u.id AS agent_id, u.first_name AS agent_first_name,
			   u.last_name AS agent_last_name, u.email AS agent_email
		FROM new_listing nl
		INNER JOIN users u ON nl.agent_id = u.id
	`

	newListing := *listing
	newListing.Agent = new(domain.Agent)

	err = tx.QueryRow(ctx, listingQuery, listing.Address, listing.Price, listing.Beds, listing.Baths, listing.SqFt, listing.AgentID).
		Scan(
			&newListing.ID,
			&newListing.CreatedAt,
			&newListing.UpdatedAt,
			&newListing.AgentID,
			&newListing.Agent.ID,
			&newListing.Agent.FirstName,
			&newListing.Agent.LastName,
			&newListing.Agent.Email,
		)
	if err != nil {
		return nil, err
	}

	imageQuery := `
		INSERT INTO listing_images (public_id, listing_id, url, is_primary, sort_order)	
		VALUES($1, $2, $3, $4, $5)
		RETURNING id, public_id, listing_id, url, is_primary, sort_order, created_at, updated_at
	`

	for _, image := range listing.Images {
		batch.Queue(
			imageQuery,
			image.PublicID,
			newListing.ID,
			image.URL,
			image.IsPrimary,
			image.SortOrder,
		)
	}

	br := tx.SendBatch(ctx, batch)

	for i := range listing.Images {
		err = br.QueryRow().Scan(
			&newListing.Images[i].ID,
			&newListing.Images[i].PublicID,
			&newListing.Images[i].ListingID,
			&newListing.Images[i].URL,
			&newListing.Images[i].IsPrimary,
			&newListing.Images[i].SortOrder,
			&newListing.Images[i].CreatedAt,
			&newListing.Images[i].UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	if err = br.Close(); err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &newListing, nil
}

func (r *ListingRepository) UpdateListingById(
	ctx context.Context,
	listing *dto.UpdateListingRequest,
	currentUserCtx *domain.ContextSessionData,
	listingId int,
) (*domain.Listing, error) {
	query := `
			UPDATE listings
			SET address = COALESCE($1, address),
				price = COALESCE($2, price),
				beds = COALESCE($3, beds),
				baths = COALESCE($4, baths),
				sq_ft = COALESCE($5, sq_ft),
				description = COALESCE($6, description),
				agent_id = COALESCE($7, agent_id),
				updated_at = NOW()
			WHERE id = $8 AND 
			(
				agent_id = $9
				OR $10 = 'admin'
			)
			RETURNING *
		`

	var updatedListing domain.Listing

	err := r.db.QueryRow(
		ctx,
		query,
		listing.Address,
		listing.Price,
		listing.Beds,
		listing.Baths,
		listing.SqFt,
		listing.Description,
		listing.AgentID,
		listingId,
		currentUserCtx.UserID,
		currentUserCtx.Role,
	).Scan(
		&updatedListing.ID,
		&updatedListing.Address,
		&updatedListing.Price,
		&updatedListing.Beds,
		&updatedListing.Baths,
		&updatedListing.SqFt,
		&updatedListing.Description,
		&updatedListing.AgentID,
		&updatedListing.CreatedAt,
		&updatedListing.UpdatedAt,
		&updatedListing.Views,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("Listing not found or you do not have permission")
		}
		return nil, err
	}
	return &updatedListing, err
}

func (r *ListingRepository) DeleteListingById(
	ctx context.Context,
	currentUserCtx *domain.ContextSessionData,
	listingId int,
) error {
	query := `
		DELETE FROM listings
		WHERE id = $1 AND 
		(
			agent_id = $2
			OR $3 = 'admin'
		)
	`

	result, err := r.db.Exec(
		ctx,
		query,
		listingId,
		currentUserCtx.UserID,
		currentUserCtx.Role,
	)
	if err != nil {
		return err
	}

	rows := result.RowsAffected()

	if rows != 1 {
		return errors.New("Cannot delete other agent's listing")
	}

	return nil
}

func (r *ListingRepository) GetAgentIdByListingId(ctx context.Context, listingId int) (int, error) {
	query := `
		SELECT agent_id FROM listings
		WHERE id = $1
	`

	var agentId int

	if err := r.db.QueryRow(ctx, query, listingId).Scan(&agentId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errors.New("Listing not found or you do not have permission")
		}
		return 0, err
	}

	return agentId, nil
}

func (r *ListingRepository) TrackViewsByListingId(
	ctx context.Context,
	listingId int,
) error {
	query := `
		UPDATE listings
		SET views = views + 1
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, listingId)
	if err != nil {
		return err
	}

	return nil
}

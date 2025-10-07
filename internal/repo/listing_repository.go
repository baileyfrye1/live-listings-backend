package repo

import (
	"context"
	"database/sql"
	"errors"

	"server/internal/api/dto"
	"server/internal/domain"
)

type IListingRepo interface {
	GetAllListings(ctx context.Context) ([]*domain.Listing, error)
	GetListingById(ctx context.Context, id int) (*domain.Listing, error)
	GetListingsByAgentId(ctx context.Context, agentId int) ([]*domain.Listing, error)
	CreateListing(ctx context.Context, listing *dto.CreateListingRequest) (*domain.Listing, error)
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
}

type ListingRepository struct {
	db *sql.DB
}

func NewListingRepository(db *sql.DB) *ListingRepository {
	return &ListingRepository{db: db}
}

func (r *ListingRepository) GetAllListings(ctx context.Context) ([]*domain.Listing, error) {
	query := `
		SELECT listings.*,
			users.id,
			users.first_name,
			users.last_name,
			users.email
		FROM listings
		INNER JOIN users
			ON listings.agent_id = users.id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var listings []*domain.Listing
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
			&listing.Agent.ID,
			&listing.Agent.FirstName,
			&listing.Agent.LastName,
			&listing.Agent.Email,
		)
		if err != nil {
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
		SELECT listings.*,
			users.id,
			users.first_name,
			users.last_name,
			users.email
		FROM listings
		INNER JOIN users
			ON listings.agent_id = users.id
		WHERE listings.id = $1
	`

	var listing domain.Listing
	listing.Agent = new(domain.Agent)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
		&listing.Agent.ID,
		&listing.Agent.FirstName,
		&listing.Agent.LastName,
		&listing.Agent.Email,
	)
	if err != nil {
		return nil, err
	}

	return &listing, nil
}

func (r *ListingRepository) GetListingsByAgentId(
	ctx context.Context,
	agentId int,
) ([]*domain.Listing, error) {
	query := `
		SELECT * FROM listings
		WHERE agent_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, agentId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var listings []*domain.Listing
	for rows.Next() {
		listing := new(domain.Listing)

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
		)
		if err != nil {
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
	listing *dto.CreateListingRequest,
) (*domain.Listing, error) {
	query := `
		WITH new_listing AS (
			INSERT INTO listings (address, price, beds, baths, sq_ft, agent_id)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING *
		)
		SELECT nl.*, u.id AS agent_id, u.first_name AS agent_first_name,
			   u.last_name AS agent_last_name, u.email AS agent_email
		FROM new_listing nl
		INNER JOIN users u ON nl.agent_id = u.id
	`

	var newListing domain.Listing
	newListing.Agent = new(domain.Agent)

	err := r.db.QueryRowContext(ctx, query, listing.Address, listing.Price, listing.Beds, listing.Baths, listing.SqFt, listing.AgentID).
		Scan(
			&newListing.ID,
			&newListing.Address,
			&newListing.Price,
			&newListing.Beds,
			&newListing.Baths,
			&newListing.SqFt,
			&newListing.Description,
			&newListing.AgentID,
			&newListing.CreatedAt,
			&newListing.UpdatedAt,
			&newListing.Agent.ID,
			&newListing.Agent.FirstName,
			&newListing.Agent.LastName,
			&newListing.Agent.Email,
		)
	if err != nil {
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

	err := r.db.QueryRowContext(
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
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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

	result, err := r.db.ExecContext(
		ctx,
		query,
		listingId,
		currentUserCtx.UserID,
		currentUserCtx.Role,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return errors.New("Cannot delete other agent's listing")
	}

	return nil
}

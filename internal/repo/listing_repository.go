package repo

import (
	"context"
	"database/sql"

	"server/internal/api/dto"
	"server/internal/domain"
)

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
		listing.Agent = new(domain.ListingAgent)

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
	listing.Agent = new(domain.ListingAgent)

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
	listing *dto.RequestCreateListing,
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
	newListing.Agent = new(domain.ListingAgent)

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

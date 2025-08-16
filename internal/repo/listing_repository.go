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
		SELECT * FROM listings
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

		err := rows.Scan(
			listing.ID,
			listing.Address,
			listing.Price,
			listing.Beds,
			listing.Baths,
			listing.SqFt,
			listing.AgentID,
			listing.CreatedAt,
			listing.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		listings = append(listings, listing)
	}

	return listings, nil
}

func (r *ListingRepository) GetListingById(ctx context.Context, id int) (*domain.Listing, error) {
	query := `
		SELECT * FROM listings
		WHERE id = $1
		INNER JOIN users
			ON listings.agent_id = users.id
	`

	var listing domain.Listing

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&listing.ID,
		&listing.Address,
		&listing.Price,
		&listing.Beds,
		&listing.Baths,
		&listing.SqFt,
		&listing.AgentID,
		&listing.CreatedAt,
		&listing.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &listing, nil
}

func (r *ListingRepository) CreateListing(
	ctx context.Context,
	listing *dto.RequestCreateListing,
) (*domain.Listing, error) {
	query := `
		INSERT INTO listings
		(address, price, beds, baths, sq_ft, agent_id)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING *
	`

	var newListing domain.Listing

	err := r.db.QueryRowContext(ctx, query).Scan(
		&newListing.ID,
		&newListing.Address,
		&newListing.Price,
		&newListing.Beds,
		&newListing.Baths,
		&newListing.SqFt,
		&newListing.AgentID,
		&newListing.CreatedAt,
		&newListing.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &newListing, nil
}

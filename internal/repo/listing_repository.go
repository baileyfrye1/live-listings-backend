package repo

import (
	"context"
	"database/sql"

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

func (r *ListingRepository) CreateListing(
	ctx context.Context,
	listing *domain.Listing,
) (*domain.Listing, error) {
	query := ``

	return nil, nil
}

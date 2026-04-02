package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"server/internal/api/dto"
	"server/internal/domain"
)

type IUserRepo interface {
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	GetUserById(ctx context.Context, id int) (*domain.User, error)
	GetAgentById(ctx context.Context, id int) (*domain.Agent, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAllAgents(ctx context.Context) ([]*domain.Agent, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateUserById(
		ctx context.Context,
		userReq *dto.UpdateUserRequest,
		id int,
	) (*domain.User, error)
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, role
		FROM users
		WHERE role = 'user' OR role = 'agent'
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*domain.User

	for rows.Next() {
		user := new(domain.User)

		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Role,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id int) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, created_at, updated_at, role
		FROM users
		WHERE id = $1
	`
	var user domain.User

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("Query user by id: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetAgentById(ctx context.Context, id int) (*domain.Agent, error) {
	query := `
		SELECT id, first_name, last_name, email, created_at, updated_at
		FROM users
		WHERE id = $1 AND role = 'agent'
	`
	var agent domain.Agent

	err := r.db.QueryRow(ctx, query, id).Scan(
		&agent.ID,
		&agent.FirstName,
		&agent.LastName,
		&agent.Email,
		&agent.CreatedAt,
		&agent.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("Query user by id: %w", err)
	}

	return &agent, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT * FROM users
		WHERE email = $1
	`

	var user domain.User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("Query user by email: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetAllAgents(ctx context.Context) ([]*domain.Agent, error) {
	query := `
        SELECT id, first_name, last_name, email, created_at, updated_at
        FROM users
        WHERE role = 'agent'
    `

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var agents []*domain.Agent
	for rows.Next() {
		agent := new(domain.Agent)

		err := rows.Scan(
			&agent.ID,
			&agent.FirstName,
			&agent.LastName,
			&agent.Email,
			&agent.CreatedAt,
			&agent.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		agents = append(agents, agent)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	listingsQuery := `
        SELECT l.*,
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
        LEFT JOIN listing_images li
            ON l.id = li.listing_id
        WHERE l.agent_id = $1
        GROUP BY l.id
    `

	for _, agent := range agents {
		listingRows, err := r.db.Query(ctx, listingsQuery, agent.ID)
		if err != nil {
			return nil, err
		}

		var listings []domain.Listing
		var imagesJson []byte
		for listingRows.Next() {
			listing := domain.Listing{}
			err := listingRows.Scan(
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
				&imagesJson,
			)
			if err != nil {
				listingRows.Close()
				return nil, err
			}

			if err = json.Unmarshal(imagesJson, &listing.Images); err != nil {
				listingRows.Close()
				return nil, err
			}

			listings = append(listings, listing)
		}

		if err = listingRows.Err(); err != nil {
			listingRows.Close()
			return nil, err
		}

		listingRows.Close()
		agent.Listings = listings
	}

	return agents, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		INSERT into users (first_name, last_name, email, password_hash, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	newUser := *user

	err := r.db.QueryRow(ctx, query, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.Role).
		Scan(&newUser.ID, &newUser.CreatedAt, &newUser.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("Insert user: %w", err)
	}

	return &newUser, nil
}

func (r *UserRepository) UpdateUserById(
	ctx context.Context,
	user *dto.UpdateUserRequest,
	id int,
) (*domain.User, error) {
	query := `
		UPDATE users
		SET first_name = COALESCE($1, first_name),
			last_name  = COALESCE($2, last_name),
			email      = COALESCE($3, email),
			role       = COALESCE($4, role),
			updated_at = NOW()
		WHERE id = $5
		RETURNING id, first_name, last_name, email, created_at, updated_at, role
	`

	var updatedUser domain.User

	err := r.db.QueryRow(ctx, query, user.FirstName, user.LastName, user.Email, user.Role, id).
		Scan(
			&updatedUser.ID,
			&updatedUser.FirstName,
			&updatedUser.LastName,
			&updatedUser.Email,
			&updatedUser.CreatedAt,
			&updatedUser.UpdatedAt,
			&updatedUser.Role,
		)
	if err != nil {
		return nil, fmt.Errorf("Update user: %w", err)
	}

	return &updatedUser, nil
}

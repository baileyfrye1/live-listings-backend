package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"server/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByUserId(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, username, email, created_at, updated_at
		FROM users
		Where id = $1
	`
	var user domain.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Query user by id: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		INSERT into users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	newUser := *user

	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.PasswordHash).
		Scan(&newUser.ID, &newUser.CreatedAt, &newUser.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("Insert user: %w", err)
	}

	return &newUser, nil
}

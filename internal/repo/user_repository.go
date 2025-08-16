package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"server/internal/api/dto"
	"server/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserById(ctx context.Context, id int) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, created_at, updated_at, role
		FROM users
		WHERE id = $1
	`
	var user domain.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("Query user by id: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT * FROM users
		WHERE email = $1
	`

	var user domain.User

	err := r.db.QueryRowContext(ctx, query, email).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("Query user by email: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUsersByRole(ctx context.Context, role string) ([]*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, created_at, updated_at, role
		FROM users
		WHERE role = $1
	`

	rows, err := r.db.QueryContext(ctx, query, role)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := new(domain.User)

		err := rows.Scan(
			user.ID,
			user.FirstName,
			user.LastName,
			user.Email,
			user.CreatedAt,
			user.UpdatedAt,
			user.Role,
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

func (r *UserRepository) GetAuthorizedUser(ctx context.Context, id int) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, created_at, updated_at, role
		FROM users
		WHERE (role = $1 OR role = $2) AND id = $3
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, "agent", "admin", id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("Query authorized user: %w", err)
	}

	fmt.Println(user)

	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		INSERT into users (first_name, last_name, email, password_hash, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	newUser := *user

	err := r.db.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.Role).
		Scan(&newUser.ID, &newUser.CreatedAt, &newUser.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("Insert user: %w", err)
	}

	return &newUser, nil
}

func (r *UserRepository) UpdateUserById(
	ctx context.Context,
	user *dto.RequestUpdateUser,
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

	err := r.db.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Email, user.Role, id).
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

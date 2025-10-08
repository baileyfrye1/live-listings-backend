package repo

import (
	"context"

	"server/internal/api/dto"
	"server/internal/domain"
)

type UserRepoMock struct {
	GetUserByIdFunc    func(ctx context.Context, id int) (*domain.User, error)
	GetAgentByIdFunc   func(ctx context.Context, id int) (*domain.Agent, error)
	GetUserByEmailFunc func(ctx context.Context, email string) (*domain.User, error)
	GetUsersByRoleFunc func(ctx context.Context, role string) ([]*domain.User, error)
	CreateUserFunc     func(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateUserByIdFunc func(ctx context.Context, user *dto.UpdateUserRequest, id int) (*domain.User, error)
}

func (u *UserRepoMock) GetUserById(ctx context.Context, id int) (*domain.User, error) {
	return u.GetUserByIdFunc(ctx, id)
}

func (u *UserRepoMock) GetAgentById(ctx context.Context, id int) (*domain.Agent, error) {
	return u.GetAgentByIdFunc(ctx, id)
}

func (u *UserRepoMock) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return u.GetUserByEmailFunc(ctx, email)
}

func (u *UserRepoMock) GetUsersByRole(ctx context.Context, role string) ([]*domain.User, error) {
	return u.GetUsersByRole(ctx, role)
}

func (u *UserRepoMock) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return u.CreateUserFunc(ctx, user)
}

func (u *UserRepoMock) UpdateUserById(
	ctx context.Context,
	user *dto.UpdateUserRequest,
	id int,
) (*domain.User, error) {
	return u.UpdateUserByIdFunc(ctx, user, id)
}

package service

import (
	"context"

	"server/internal/api/dto"
	"server/internal/domain"
	"server/internal/repo"
)

type UserService struct {
	userRepo *repo.UserRepository
}

func NewUserService(userRepo *repo.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserById(ctx context.Context, id int) (*domain.User, error) {
	return s.userRepo.GetUserById(ctx, id)
}

func (s *UserService) GetUsersByRole(ctx context.Context, role string) ([]*domain.User, error) {
	return s.userRepo.GetUsersByRole(ctx, role)
}

func (s *UserService) UpdateUserById(
	ctx context.Context,
	user *dto.RequestUpdateUser,
	id int,
) (*domain.User, error) {
	return s.userRepo.UpdateUserById(ctx, user, id)
}

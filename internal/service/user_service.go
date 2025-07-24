package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/domain"
	"server/internal/repo"
)

type UserService struct {
	userRepo *repo.UserRepository
}

func NewUserService(userRepo *repo.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByUserId(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context) {
}

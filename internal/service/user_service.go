package service

import (
	"context"
	"errors"

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

func (s *UserService) GetAgentById(ctx context.Context, id int) (*domain.Agent, error) {
	return s.userRepo.GetAgentById(ctx, id)
}

func (s *UserService) GetUsersByRole(ctx context.Context, role string) ([]*domain.User, error) {
	return s.userRepo.GetUsersByRole(ctx, role)
}

func (s *UserService) UpdateUserById(
	ctx context.Context,
	userReq *dto.UpdateUserRequest,
	userCtx *domain.ContextSessionData,
) (*domain.User, error) {
	if userCtx.Role != "admin" && userReq.Role != nil && *userReq.Role == "admin" {
		return nil, errors.New(
			"Cannot change role to admin. Please contact admin to request admin privilages",
		)
	}

	return s.userRepo.UpdateUserById(ctx, userReq, userCtx.UserID)
}

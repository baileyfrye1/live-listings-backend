package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"server/internal/api/dto"
	"server/internal/domain"
	"server/internal/repo"
	"server/internal/session"
	"server/util"
)

type AuthService struct {
	userRepo *repo.UserRepository
	session  *session.Session
}

func NewAuthService(userRepo *repo.UserRepository, session *session.Session) *AuthService {
	return &AuthService{userRepo: userRepo, session: session}
}

func (s *AuthService) Register(
	ctx context.Context,
	req *dto.CreateUserRequest,
) (*dto.LoginUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("Error creating user: %v\n", err)
	}

	if req.Role == "" {
		req.Role = "user"
	}

	u := &domain.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Role:         req.Role,
		PasswordHash: &hashedPassword,
	}

	newUser, err := s.userRepo.CreateUser(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user: %v\n", err)
	}

	sessionId, err := util.CreateSession(ctx, s.session, newUser)
	if err != nil {
		return nil, err
	}

	return &dto.LoginUserResponse{
		ID:        newUser.ID,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email,
		Role:      newUser.Role,
		SessionID: sessionId,
	}, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	req *dto.LoginUserRequest,
) (*dto.LoginUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("User does not exist. Please create an account")
	}

	if err = util.CompareHashedPassword(*user.PasswordHash, req.Password); err != nil {
		return nil, errors.New("Invalid username/password")
	}

	sessionId, err := util.CreateSession(ctx, s.session, user)
	if err != nil {
		return nil, err
	}

	return &dto.LoginUserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		SessionID: sessionId,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionId string) error {
	return s.session.DeleteSession(ctx, sessionId)
}

package service

import (
	"context"
	"testing"
	"time"

	"server/internal/api/dto"
	"server/internal/domain"
	"server/internal/repo"
	"server/internal/session"
	"server/util"
)

func TestRegister(t *testing.T) {
	mockRepo := newMockUserRepo()
	mockSession := newMockSession()

	t.Run("Valid registration information returns success", func(t *testing.T) {
		ctx := context.Background()
		req := &dto.CreateUserRequest{
			FirstName: "Bailey",
			LastName:  "Frye",
			Email:     "test@test.com",
			Password:  "password",
			Role:      "agent",
		}

		a := NewAuthService(mockRepo, mockSession)

		res, err := a.Register(ctx, req)
		if err != nil {
			t.Errorf("Expected success, received %v", err.Error())
		}

		if res.FirstName != req.FirstName {
			t.Errorf("Expected %s, received %s", req.FirstName, res.FirstName)
		}

		if res.LastName != req.LastName {
			t.Errorf("Expected %s, received %s", req.LastName, res.LastName)
		}

		if res.Email != req.Email {
			t.Errorf("Expected %s, received %s", req.Email, res.Email)
		}

		if res.Role != req.Role {
			t.Errorf("Expected %s, received %s", req.Role, res.Role)
		}
	})

	t.Run("User request omitting role returns success", func(t *testing.T) {
		mockRepo := &repo.UserRepoMock{
			CreateUserFunc: func(ctx context.Context, user *domain.User) (*domain.User, error) {
				return &domain.User{
					ID:        1,
					FirstName: "Bailey",
					LastName:  "Frye",
					Email:     "test@test.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Role:      "user",
				}, nil
			},
		}

		ctx := context.Background()
		req := &dto.CreateUserRequest{
			FirstName: "Bailey",
			LastName:  "Frye",
			Email:     "test@test.com",
			Password:  "password",
		}

		a := NewAuthService(mockRepo, mockSession)

		res, err := a.Register(ctx, req)
		if err != nil {
			t.Errorf("Expected success, received %v", err.Error())
		}

		if res.FirstName != req.FirstName {
			t.Errorf("Expected %s, received %s", req.FirstName, res.FirstName)
		}

		if res.LastName != req.LastName {
			t.Errorf("Expected %s, received %s", req.LastName, res.LastName)
		}

		if res.Email != req.Email {
			t.Errorf("Expected %s, received %s", req.Email, res.Email)
		}

		if req.Role != "user" {
			t.Errorf("Expected request role to be updated to user, received %s", req.Role)
		}

		if res.Role != "user" {
			t.Errorf("Expected user, received %s", res.Role)
		}
	})

	validationTests := []struct {
		Name    string
		UserReq *dto.CreateUserRequest
	}{
		{
			Name: "User request missing first name returns error",
			UserReq: &dto.CreateUserRequest{
				LastName: "Frye",
				Email:    "test@test.com",
				Password: "password",
				Role:     "agent",
			},
		},
		{
			Name: "User request missing last name returns error",
			UserReq: &dto.CreateUserRequest{
				FirstName: "Bailey",
				Email:     "test@test.com",
				Password:  "password",
				Role:      "agent",
			},
		},
		{
			Name: "User request missing email returns error",
			UserReq: &dto.CreateUserRequest{
				FirstName: "Bailey",
				LastName:  "Frye",
				Password:  "password",
				Role:      "agent",
			},
		},
		{
			Name: "User request missing password returns error",
			UserReq: &dto.CreateUserRequest{
				FirstName: "Bailey",
				LastName:  "Frye",
				Email:     "test@test.com",
				Role:      "agent",
			},
		},
	}

	for _, vt := range validationTests {
		t.Run(vt.Name, func(t *testing.T) {
			ctx := context.Background()
			req := vt.UserReq

			a := NewAuthService(mockRepo, mockSession)

			_, err := a.Register(ctx, req)
			wantErr := "Please enter all fields"

			if err == nil {
				t.Errorf("Expected %s, received nil", wantErr)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	mockRepo := newMockUserRepo()
	mockSession := newMockSession()

	loginTests := []struct {
		Name      string
		LoginReq  *dto.LoginUserRequest
		ExpectErr bool
	}{
		{
			Name: "Valid login credentials returns success",
			LoginReq: &dto.LoginUserRequest{
				Email:    "test@test.com",
				Password: "password",
			},
			ExpectErr: false,
		},
		{
			Name: "Missing email returns error",
			LoginReq: &dto.LoginUserRequest{
				Password: "password",
			},
			ExpectErr: true,
		},
		{
			Name: "Missing password returns error",
			LoginReq: &dto.LoginUserRequest{
				Email: "test@test.com",
			},
			ExpectErr: true,
		},
		{
			Name: "Incorrect password returns error",
			LoginReq: &dto.LoginUserRequest{
				Email:    "test@test.com",
				Password: "password123",
			},
			ExpectErr: true,
		},
	}

	for _, lt := range loginTests {
		t.Run(lt.Name, func(t *testing.T) {
			ctx := context.Background()
			req := lt.LoginReq

			a := NewAuthService(mockRepo, mockSession)

			_, err := a.Login(ctx, req)

			if lt.ExpectErr && err == nil {
				t.Error("Expected error, received nil")
			}

			if !lt.ExpectErr && err != nil {
				t.Errorf("Expected success, received %v", err)
			}
		})
	}
}

func newMockUserRepo() *repo.UserRepoMock {
	hashedPassword, _ := util.HashPassword("password")
	return &repo.UserRepoMock{
		GetUserByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{
				ID:           1,
				FirstName:    "Bailey",
				LastName:     "Frye",
				Email:        "test@test.com",
				PasswordHash: hashedPassword,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				Role:         "agent",
			}, nil
		},
		CreateUserFunc: func(ctx context.Context, user *domain.User) (*domain.User, error) {
			return &domain.User{
				ID:        1,
				FirstName: "Bailey",
				LastName:  "Frye",
				Email:     "test@test.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Role:      "agent",
			}, nil
		},
	}
}

func newMockSession() *session.SessionMock {
	return &session.SessionMock{
		GetSessionFunc: func(ctx context.Context, sessionID string) (*domain.SessionData, error) {
			return nil, nil
		},
		SetSessionFunc: func(ctx context.Context, sessionID string, sessionData *domain.SessionData, ttl time.Duration) error {
			return nil
		},
		DeleteSessionFunc: func(ctx context.Context, sessionID string) error {
			return nil
		},
	}
}

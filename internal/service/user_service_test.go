package service

import (
	"context"
	"testing"

	"server/internal/api/dto"
	"server/internal/domain"
	"server/internal/repo"
)

func TestUpdateUser(t *testing.T) {
	t.Run("Agent tries to change role to admin returns error", func(t *testing.T) {
		mockRepo := &repo.UserRepoMock{
			UpdateUserByIdFunc: func(ctx context.Context, user *dto.UpdateUserRequest, id int) (*domain.User, error) {
				return &domain.User{ID: 123, Role: "agent"}, nil
			},
		}

		reqRole := "admin"
		userReq := &dto.UpdateUserRequest{Role: &reqRole}

		userCtx := &domain.ContextSessionData{SessionID: "123abc", UserID: 123, Role: "agent"}

		ctx := context.Background()

		u := NewUserService(mockRepo)
		_, err := u.UpdateUserById(ctx, userReq, userCtx)
		wantErr := "Cannot change role to admin. Please contact admin to request admin privileges"

		if err == nil {
			t.Errorf("Expected err %q, received nil", wantErr)
		}

		if err.Error() != wantErr {
			t.Errorf("Got %q want %q", err.Error(), wantErr)
		}
	})

	t.Run("User tries to change role to admin returns error", func(t *testing.T) {
		mockRepo := &repo.UserRepoMock{
			UpdateUserByIdFunc: func(ctx context.Context, user *dto.UpdateUserRequest, id int) (*domain.User, error) {
				return &domain.User{ID: 123, Role: "user"}, nil
			},
		}

		reqRole := "admin"
		userReq := &dto.UpdateUserRequest{Role: &reqRole}

		userCtx := &domain.ContextSessionData{SessionID: "123abc", UserID: 123, Role: "user"}

		ctx := context.Background()

		u := NewUserService(mockRepo)
		_, err := u.UpdateUserById(ctx, userReq, userCtx)
		wantErr := "Cannot change role to admin. Please contact admin to request admin privileges"

		if err == nil {
			t.Errorf("Expected err %q, received nil", wantErr)
		}

		if err.Error() != wantErr {
			t.Errorf("Got %q want %q", err.Error(), wantErr)
		}
	})

	t.Run("Admin tries to change role to admin returns success", func(t *testing.T) {
		mockRepo := &repo.UserRepoMock{
			UpdateUserByIdFunc: func(ctx context.Context, user *dto.UpdateUserRequest, id int) (*domain.User, error) {
				return &domain.User{ID: 123, Role: "admin"}, nil
			},
		}

		reqRole := "admin"
		userReq := &dto.UpdateUserRequest{Role: &reqRole}

		userCtx := &domain.ContextSessionData{SessionID: "123abc", UserID: 123, Role: "admin"}

		ctx := context.Background()

		u := NewUserService(mockRepo)
		_, err := u.UpdateUserById(ctx, userReq, userCtx)
		if err != nil {
			t.Errorf("Expected success, received %q", err.Error())
		}
	})
}

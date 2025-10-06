package service

import (
	"context"
	"testing"

	"server/internal/api/dto"
	"server/internal/domain"
)

type userRepoMock struct {
	UpdateUserByIdFunc func(ctx context.Context, user *dto.UpdateUserRequest, id int) (*domain.User, error)
}

func (u *userRepoMock) UpdateUserById(
	ctx context.Context,
	user *dto.UpdateUserRequest,
	id int,
) (*domain.User, error) {
	return u.UpdateUserByIdFunc(ctx, user, id)
}

func TestUpdateUser(t *testing.T) {
	t.Run("Agent tries to change role to admin returns error", func(t *testing.T) {
		mockRepo := &userRepoMock{
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
		mockRepo := &userRepoMock{
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
		mockRepo := &userRepoMock{
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

// Rest of interface methods to satisfy interface requirements
func (u *userRepoMock) GetUserById(ctx context.Context, id int) (*domain.User, error) {
	return nil, nil
}

func (u *userRepoMock) GetAgentById(ctx context.Context, id int) (*domain.Agent, error) {
	return nil, nil
}

func (u *userRepoMock) GetUsersByRole(ctx context.Context, role string) ([]*domain.User, error) {
	return nil, nil
}

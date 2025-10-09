package service

import (
	"context"
	"testing"

	"server/internal/api/dto"
	"server/internal/domain"
	"server/internal/repo"
)

func TestGetAllUsers(t *testing.T) {
	tests := []struct {
		Name      string
		UserCtx   *domain.ContextSessionData
		ExpectErr bool
	}{
		{
			Name: "Agent tries to view all users returns error",
			UserCtx: &domain.ContextSessionData{
				SessionID: "abc123",
				UserID:    1,
				Role:      "agent",
			},
			ExpectErr: true,
		},
		{
			Name: "Admin tries to view all users returns success",
			UserCtx: &domain.ContextSessionData{
				SessionID: "abc123",
				UserID:    1,
				Role:      "admin",
			},
			ExpectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockRepo := &repo.UserRepoMock{
				GetAllUsersFunc: func(ctx context.Context) ([]*domain.User, error) {
					return []*domain.User{
						{
							ID:        1,
							FirstName: "Bailey",
							LastName:  "Frye",
							Email:     "test@test.com",
						},

						{
							ID:        2,
							FirstName: "Arlo",
							LastName:  "Jenkins",
							Email:     "arlo@test.com",
						},
					}, nil
				},
			}

			ctx := context.Background()
			userCtx := tt.UserCtx

			u := NewUserService(mockRepo)
			_, err := u.GetAllUsers(ctx, userCtx)

			if tt.ExpectErr && err == nil {
				t.Error("Expected err, received nil")
			}

			if !tt.ExpectErr && err != nil {
				t.Errorf("Expected success, received %v", err)
			}
		})
	}
}

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
		_, err := u.UpdateUserById(ctx, userReq, userCtx, userCtx.UserID)
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
		_, err := u.UpdateUserById(ctx, userReq, userCtx, userCtx.UserID)
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
		_, err := u.UpdateUserById(ctx, userReq, userCtx, userCtx.UserID)
		if err != nil {
			t.Errorf("Expected success, received %q", err.Error())
		}
	})
}

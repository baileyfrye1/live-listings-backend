package service

import (
	"context"
	"reflect"
	"testing"
	"time"

	"server/internal/domain"
	"server/internal/repo"
)

func TestGetAllFavorites(t *testing.T) {
	fixedTime := time.Date(2025, 10, 13, 8, 0, 0, 0, time.UTC)
	tests := []struct {
		Name           string
		MockRepo       *repo.FavoriteRepoMock
		ExpectedResult []*domain.Favorite
	}{
		{
			Name: "List of favorites returned from database returns list of favorites",
			MockRepo: &repo.FavoriteRepoMock{
				GetUserFavoritesFunc: func(ctx context.Context, userCtx *domain.ContextSessionData) ([]*domain.Favorite, error) {
					return []*domain.Favorite{
						{
							ID:        1,
							UserID:    1,
							ListingID: 1,
							CreatedAt: fixedTime,
							UpdatedAt: fixedTime,
						},
						{
							ID:        2,
							UserID:    2,
							ListingID: 2,
							CreatedAt: fixedTime,
							UpdatedAt: fixedTime,
						},
					}, nil
				},
			},
			ExpectedResult: []*domain.Favorite{
				{
					ID:        1,
					UserID:    1,
					ListingID: 1,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				{
					ID:        2,
					UserID:    2,
					ListingID: 2,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
		},
		{
			Name: "No favorites in database returns empty slice",
			MockRepo: &repo.FavoriteRepoMock{
				GetUserFavoritesFunc: func(ctx context.Context, userCtx *domain.ContextSessionData) ([]*domain.Favorite, error) {
					return nil, nil
				},
			},
			ExpectedResult: []*domain.Favorite{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockRepo := tt.MockRepo
			ctx := context.Background()
			userCtx := &domain.ContextSessionData{
				SessionID: "abc123",
				UserID:    1,
				Role:      "user",
			}

			f := NewFavoriteService(mockRepo)

			favorites, err := f.GetUserFavorites(ctx, userCtx)
			if err != nil {
				t.Error("Wanted empty slice or slice of favorites, received error")
			}

			if !reflect.DeepEqual(favorites, tt.ExpectedResult) {
				t.Errorf("Expected %#v, received %#v", favorites, tt.ExpectedResult)
			}
		})
	}
}

func TestGetFavoritesMap(t *testing.T) {
	tests := []struct {
		Name           string
		MockRepo       *repo.FavoriteRepoMock
		ExpectedResult map[int]bool
	}{
		{
			Name: "List of favorites from repo returns valid map",
			MockRepo: &repo.FavoriteRepoMock{
				GetUserFavoritesFunc: func(ctx context.Context, userCtx *domain.ContextSessionData) ([]*domain.Favorite, error) {
					return []*domain.Favorite{
						{
							ID:        1,
							UserID:    4,
							ListingID: 2,
						},
						{
							ID:        2,
							UserID:    4,
							ListingID: 6,
						},
						{
							ID:        3,
							UserID:    4,
							ListingID: 30,
						},
					}, nil
				},
			},
			ExpectedResult: map[int]bool{
				2:  true,
				6:  true,
				30: true,
			},
		},
		{
			Name: "No favorites in database returns empty map",
			MockRepo: &repo.FavoriteRepoMock{
				GetUserFavoritesFunc: func(ctx context.Context, userCtx *domain.ContextSessionData) ([]*domain.Favorite, error) {
					return nil, nil
				},
			},
			ExpectedResult: map[int]bool{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockRepo := tt.MockRepo
			ctx := context.Background()
			userCtx := &domain.ContextSessionData{
				SessionID: "abc123",
				UserID:    1,
				Role:      "user",
			}

			f := NewFavoriteService(mockRepo)

			favorites, err := f.GetUserFavoritesMap(ctx, userCtx)
			if err != nil {
				t.Error("Wanted empty slice or slice of favorites, received error")
			}

			if !reflect.DeepEqual(favorites, tt.ExpectedResult) {
				t.Errorf("Expected %#v, received %#v", favorites, tt.ExpectedResult)
			}
		})
	}
}

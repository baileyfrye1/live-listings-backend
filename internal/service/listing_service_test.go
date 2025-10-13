package service

import (
	"context"
	"reflect"
	"testing"
	"time"

	"server/internal/api/dto"
	"server/internal/domain"
	"server/internal/repo"
)

func TestGetAllListings(t *testing.T) {
	tests := []struct {
		Name           string
		MockRepo       *repo.ListingRepoMock
		ExpectedResult []*domain.Listing
	}{
		{
			Name: "List of listings in database returns slice of listings",
			MockRepo: &repo.ListingRepoMock{
				GetAllListingsFunc: func(ctx context.Context) ([]*domain.Listing, error) {
					return []*domain.Listing{
						{
							ID:      1,
							Price:   400000,
							Beds:    2,
							Baths:   1,
							Address: "123 Test St, Nashville, TN",
							SqFt:    3200,
						},
						{
							ID:      2,
							Price:   300000,
							Beds:    3,
							Baths:   2,
							Address: "124 Test St, Nashville, TN",
							SqFt:    1600,
						},
						{
							ID:      3,
							Price:   500000,
							Beds:    3,
							Baths:   2,
							Address: "125 Test St, Nashville, TN",
							SqFt:    4000,
						},
					}, nil
				},
			},
			ExpectedResult: []*domain.Listing{
				{
					ID:      1,
					Price:   400000,
					Beds:    2,
					Baths:   1,
					Address: "123 Test St, Nashville, TN",
					SqFt:    3200,
				},
				{
					ID:      2,
					Price:   300000,
					Beds:    3,
					Baths:   2,
					Address: "124 Test St, Nashville, TN",
					SqFt:    1600,
				},
				{
					ID:      3,
					Price:   500000,
					Beds:    3,
					Baths:   2,
					Address: "125 Test St, Nashville, TN",
					SqFt:    4000,
				},
			},
		},
		{
			Name: "No listings in database returns empty slice",
			MockRepo: &repo.ListingRepoMock{
				GetAllListingsFunc: func(ctx context.Context) ([]*domain.Listing, error) {
					return nil, nil
				},
			},
			ExpectedResult: []*domain.Listing{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockRepo := tt.MockRepo
			ctx := context.Background()

			l := NewListingService(mockRepo)

			listings, err := l.GetAllListings(ctx)
			if err != nil {
				t.Error("Wanted empty slice or slice of favorites, received error")
			}

			if !reflect.DeepEqual(listings, tt.ExpectedResult) {
				t.Errorf("Expected %#v, received %#v", listings, tt.ExpectedResult)
			}
		})
	}
}

func TestGetAllListingsByAgentId(t *testing.T) {
	tests := []struct {
		Name           string
		MockRepo       *repo.ListingRepoMock
		ExpectedResult []*domain.Listing
	}{
		{
			Name: "List of agent listings in database returns slice of listings",
			MockRepo: &repo.ListingRepoMock{
				GetListingsByAgentIdFunc: func(ctx context.Context, agentId int) ([]*domain.Listing, error) {
					return []*domain.Listing{
						{
							ID:      1,
							Price:   400000,
							Beds:    2,
							Baths:   1,
							Address: "123 Test St, Nashville, TN",
							SqFt:    3200,
						},
						{
							ID:      2,
							Price:   300000,
							Beds:    3,
							Baths:   2,
							Address: "124 Test St, Nashville, TN",
							SqFt:    1600,
						},
						{
							ID:      3,
							Price:   500000,
							Beds:    3,
							Baths:   2,
							Address: "125 Test St, Nashville, TN",
							SqFt:    4000,
						},
					}, nil
				},
			},
			ExpectedResult: []*domain.Listing{
				{
					ID:      1,
					Price:   400000,
					Beds:    2,
					Baths:   1,
					Address: "123 Test St, Nashville, TN",
					SqFt:    3200,
				},
				{
					ID:      2,
					Price:   300000,
					Beds:    3,
					Baths:   2,
					Address: "124 Test St, Nashville, TN",
					SqFt:    1600,
				},
				{
					ID:      3,
					Price:   500000,
					Beds:    3,
					Baths:   2,
					Address: "125 Test St, Nashville, TN",
					SqFt:    4000,
				},
			},
		},
		{
			Name: "No agent listings in database returns empty slice",
			MockRepo: &repo.ListingRepoMock{
				GetListingsByAgentIdFunc: func(ctx context.Context, agentId int) ([]*domain.Listing, error) {
					return nil, nil
				},
			},
			ExpectedResult: []*domain.Listing{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockRepo := tt.MockRepo
			ctx := context.Background()
			agentId := 1

			l := NewListingService(mockRepo)

			favorites, err := l.GetListingsByAgentId(ctx, agentId)
			if err != nil {
				t.Error("Wanted empty slice or slice of favorites, received error")
			}

			if !reflect.DeepEqual(favorites, tt.ExpectedResult) {
				t.Errorf("Expected %#v, received %#v", favorites, tt.ExpectedResult)
			}
		})
	}
}

func TestUpdateListing(t *testing.T) {
	t.Run("Agent tries to change agent id on listing returns error", func(t *testing.T) {
		mockListing := &repo.ListingRepoMock{
			UpdateListingByIdFunc: func(ctx context.Context, listingReq *dto.UpdateListingRequest, userCtx *domain.ContextSessionData, id int) (*domain.Listing, error) {
				return &domain.Listing{
					ID:        1,
					Address:   "2912 River Bend Dr, Nashville, TN 37214",
					Price:     375000,
					Beds:      3,
					Baths:     2,
					SqFt:      1300,
					AgentID:   1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
		}

		reqAgentId := 3
		listingReq := &dto.UpdateListingRequest{AgentID: &reqAgentId}
		userCtx := &domain.ContextSessionData{SessionID: "123abc", UserID: 123, Role: "agent"}
		ctx := context.Background()

		l := NewListingService(mockListing)
		_, err := l.UpdateListingById(ctx, listingReq, userCtx, 1)
		wantErr := "Cannot update agent on listing. Please contact admin to change agent"

		if err == nil {
			t.Errorf("Expected err %q, received nil", wantErr)
		}

		if err.Error() != wantErr {
			t.Errorf("Got %q want %q", err.Error(), wantErr)
		}
	})

	t.Run("Admin tries to change agent id on listing returns success", func(t *testing.T) {
		mockListing := &repo.ListingRepoMock{
			UpdateListingByIdFunc: func(ctx context.Context, listingReq *dto.UpdateListingRequest, userCtx *domain.ContextSessionData, id int) (*domain.Listing, error) {
				return &domain.Listing{
					ID:        1,
					Address:   "2912 River Bend Dr, Nashville, TN 37214",
					Price:     375000,
					Beds:      3,
					Baths:     2,
					SqFt:      1300,
					AgentID:   1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
		}

		reqAgentId := 3
		listingReq := &dto.UpdateListingRequest{AgentID: &reqAgentId}
		userCtx := &domain.ContextSessionData{SessionID: "123abc", UserID: 123, Role: "admin"}
		ctx := context.Background()

		l := NewListingService(mockListing)
		_, err := l.UpdateListingById(ctx, listingReq, userCtx, 1)
		if err != nil {
			t.Errorf("Expected success, received %q", err.Error())
		}
	})
}

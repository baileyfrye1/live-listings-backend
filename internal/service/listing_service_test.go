package service

import (
	"context"
	"testing"
	"time"

	"server/internal/api/dto"
	"server/internal/domain"
)

type listingRepoMock struct {
	UpdateListingByIdFunc func(ctx context.Context, listingReq *dto.UpdateListingRequest, userCtx *domain.ContextSessionData, id int) (*domain.Listing, error)
}

func (l *listingRepoMock) UpdateListingById(
	ctx context.Context,
	listingReq *dto.UpdateListingRequest,
	userCtx *domain.ContextSessionData,
	id int,
) (*domain.Listing, error) {
	return l.UpdateListingByIdFunc(ctx, listingReq, userCtx, id)
}

func TestUpdateListing(t *testing.T) {
	t.Run("Agent tries to change agent id on listing returns error", func(t *testing.T) {
		mockListing := &listingRepoMock{
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
		mockListing := &listingRepoMock{
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

// Rest of interface methods to satisfy interface requirements
func (l *listingRepoMock) GetAllListings(ctx context.Context) ([]*domain.Listing, error) {
	return nil, nil
}

func (l *listingRepoMock) GetListingById(ctx context.Context, id int) (*domain.Listing, error) {
	return nil, nil
}

func (l *listingRepoMock) GetListingsByAgentId(
	ctx context.Context,
	agentId int,
) ([]*domain.Listing, error) {
	return nil, nil
}

func (l *listingRepoMock) CreateListing(
	ctx context.Context,
	listing *dto.CreateListingRequest,
) (*domain.Listing, error) {
	return nil, nil
}

func (l *listingRepoMock) DeleteListingById(
	ctx context.Context,
	currentUserCtx *domain.ContextSessionData,
	listingId int,
) error {
	return nil
}

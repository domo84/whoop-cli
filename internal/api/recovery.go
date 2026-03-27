package api

import (
	"context"
	"fmt"
)

// RecoveryService defines operations for the /v2/recovery/ endpoints.
type RecoveryService interface {
	List(ctx context.Context, params ListParams) (*PaginatedResponse[Recovery], error)
	ListAll(ctx context.Context, params ListParams) ([]Recovery, error)
	Get(ctx context.Context, cycleID int) (*Recovery, error)
}

type recoveryService struct{ client *Client }

// NewRecoveryService creates a RecoveryService backed by client.
func NewRecoveryService(client *Client) RecoveryService {
	return &recoveryService{client: client}
}

func (s *recoveryService) List(ctx context.Context, params ListParams) (*PaginatedResponse[Recovery], error) {
	var out PaginatedResponse[Recovery]
	if err := s.client.get(ctx, "/v2/recovery", buildQuery(params), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *recoveryService) ListAll(ctx context.Context, params ListParams) ([]Recovery, error) {
	var all []Recovery
	for {
		page, err := s.List(ctx, params)
		if err != nil {
			return nil, err
		}
		all = append(all, page.Records...)
		if page.NextToken == "" {
			break
		}
		params.NextToken = page.NextToken
	}
	return all, nil
}

// Get retrieves a recovery record by the associated cycle ID.
func (s *recoveryService) Get(ctx context.Context, cycleID int) (*Recovery, error) {
	var out Recovery
	if err := s.client.get(ctx, fmt.Sprintf("/v2/cycle/%d/recovery", cycleID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

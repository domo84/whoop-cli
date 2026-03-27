package api

import (
	"context"
	"fmt"
)

// SleepService defines operations for the /v2/activity/sleep/ endpoints.
type SleepService interface {
	List(ctx context.Context, params ListParams) (*PaginatedResponse[Sleep], error)
	ListAll(ctx context.Context, params ListParams) ([]Sleep, error)
	Get(ctx context.Context, sleepID string) (*Sleep, error)
}

type sleepService struct{ client *Client }

// NewSleepService creates a SleepService backed by client.
func NewSleepService(client *Client) SleepService {
	return &sleepService{client: client}
}

func (s *sleepService) List(ctx context.Context, params ListParams) (*PaginatedResponse[Sleep], error) {
	var out PaginatedResponse[Sleep]
	if err := s.client.get(ctx, "/v2/activity/sleep", buildQuery(params), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *sleepService) ListAll(ctx context.Context, params ListParams) ([]Sleep, error) {
	var all []Sleep
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

func (s *sleepService) Get(ctx context.Context, sleepID string) (*Sleep, error) {
	var out Sleep
	if err := s.client.get(ctx, fmt.Sprintf("/v2/activity/sleep/%s", sleepID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

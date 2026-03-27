package api

import (
	"context"
	"fmt"
)

// CycleService defines operations for the /v2/cycle/ endpoints.
type CycleService interface {
	List(ctx context.Context, params ListParams) (*PaginatedResponse[Cycle], error)
	ListAll(ctx context.Context, params ListParams) ([]Cycle, error)
	Get(ctx context.Context, cycleID int) (*Cycle, error)
	GetSleep(ctx context.Context, cycleID int) (*Sleep, error)
}

type cycleService struct{ client *Client }

// NewCycleService creates a CycleService backed by client.
func NewCycleService(client *Client) CycleService {
	return &cycleService{client: client}
}

func (s *cycleService) List(ctx context.Context, params ListParams) (*PaginatedResponse[Cycle], error) {
	var out PaginatedResponse[Cycle]
	if err := s.client.get(ctx, "/v2/cycle", buildQuery(params), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *cycleService) ListAll(ctx context.Context, params ListParams) ([]Cycle, error) {
	var all []Cycle
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

func (s *cycleService) Get(ctx context.Context, cycleID int) (*Cycle, error) {
	var out Cycle
	if err := s.client.get(ctx, fmt.Sprintf("/v2/cycle/%d", cycleID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *cycleService) GetSleep(ctx context.Context, cycleID int) (*Sleep, error) {
	var out Sleep
	if err := s.client.get(ctx, fmt.Sprintf("/v2/cycle/%d/sleep", cycleID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

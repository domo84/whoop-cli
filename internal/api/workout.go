package api

import (
	"context"
	"fmt"
)

// WorkoutService defines operations for the /v2/activity/workout/ endpoints.
type WorkoutService interface {
	List(ctx context.Context, params ListParams) (*PaginatedResponse[Workout], error)
	ListAll(ctx context.Context, params ListParams) ([]Workout, error)
	Get(ctx context.Context, workoutID string) (*Workout, error)
}

type workoutService struct{ client *Client }

// NewWorkoutService creates a WorkoutService backed by client.
func NewWorkoutService(client *Client) WorkoutService {
	return &workoutService{client: client}
}

func (s *workoutService) List(ctx context.Context, params ListParams) (*PaginatedResponse[Workout], error) {
	var out PaginatedResponse[Workout]
	if err := s.client.get(ctx, "/v2/activity/workout", buildQuery(params), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *workoutService) ListAll(ctx context.Context, params ListParams) ([]Workout, error) {
	var all []Workout
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

func (s *workoutService) Get(ctx context.Context, workoutID string) (*Workout, error) {
	var out Workout
	if err := s.client.get(ctx, fmt.Sprintf("/v2/activity/workout/%s", workoutID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

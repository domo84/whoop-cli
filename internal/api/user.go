package api

import "context"

// UserService defines operations for the /v2/users/ endpoints.
type UserService interface {
	GetProfile(ctx context.Context) (*UserProfile, error)
	GetBodyMeasurement(ctx context.Context) (*BodyMeasurement, error)
}

type userService struct{ client *Client }

// NewUserService creates a UserService backed by client.
func NewUserService(client *Client) UserService {
	return &userService{client: client}
}

func (s *userService) GetProfile(ctx context.Context) (*UserProfile, error) {
	var out UserProfile
	if err := s.client.get(ctx, "/v2/user/profile/basic", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *userService) GetBodyMeasurement(ctx context.Context) (*BodyMeasurement, error) {
	var out BodyMeasurement
	if err := s.client.get(ctx, "/v2/user/measurement/body", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

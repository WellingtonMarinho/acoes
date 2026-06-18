package devices

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, registration Registration) (Registration, error) {
	return s.repo.Upsert(ctx, registration)
}

func (s *Service) List(ctx context.Context) ([]Registration, error) {
	return s.repo.List(ctx)
}

func (s *Service) Resolve(ctx context.Context, userID string) (string, bool, error) {
	registration, ok, err := s.repo.Resolve(ctx, userID)
	if err != nil || !ok {
		return "", ok, err
	}
	return registration.DeviceToken, true, nil
}

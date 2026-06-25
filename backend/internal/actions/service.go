package actions

import (
	"context"
	"errors"
	"strings"
	"time"
)

var ErrActionNotFound = errors.New("action not found")
var ErrInvalidAction = errors.New("invalid action")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListActions(ctx context.Context, query string) ([]Action, error) {
	return s.repo.List(ctx, query)
}

func (s *Service) GetAction(ctx context.Context, id string) (Action, error) {
	action, err := s.repo.Get(ctx, strings.TrimSpace(id))
	if err != nil {
		return Action{}, err
	}
	return action, nil
}

func (s *Service) CreateAction(ctx context.Context, symbol, name, exchange string) (Action, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	name = strings.TrimSpace(name)
	exchange = strings.TrimSpace(exchange)
	if symbol == "" || name == "" {
		return Action{}, ErrInvalidAction
	}
	now := time.Now().UTC()
	action := Action{
		ID:        "action-" + strings.ToLower(symbol),
		Symbol:    symbol,
		Name:      name,
		Exchange:  exchange,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.repo.Upsert(ctx, action)
}

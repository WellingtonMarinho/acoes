package watchlist

import (
	"context"
	"errors"
	"strings"
	"time"

	"ideacoes/backend/internal/actions"
	"ideacoes/backend/internal/alerts"
	"ideacoes/backend/internal/pricefeed"
)

var (
	ErrWatchlistItemNotFound = errors.New("watchlist item not found")
	ErrInvalidWatchlistItem  = errors.New("invalid watchlist item")
)

type ActionResolver interface {
	GetAction(ctx context.Context, id string) (actions.Action, error)
}

type AlertRepository interface {
	ListByUser(ctx context.Context, userID string) ([]alerts.Alert, error)
	DeleteByUserAndAction(ctx context.Context, userID, actionID string) (int64, error)
}

type Service struct {
	repo    Repository
	actions ActionResolver
	alerts  AlertRepository
	feed    pricefeed.Feed
}

func NewService(repo Repository, actions ActionResolver, alertsRepo AlertRepository, feed pricefeed.Feed) *Service {
	return &Service{repo: repo, actions: actions, alerts: alertsRepo, feed: feed}
}

func (s *Service) Upsert(ctx context.Context, userID, actionID string) error {
	_, err := s.Add(ctx, userID, actionID)
	return err
}

func (s *Service) Add(ctx context.Context, userID, actionID string) (Item, error) {
	userID = strings.TrimSpace(userID)
	actionID = strings.TrimSpace(actionID)
	if userID == "" || actionID == "" {
		return Item{}, ErrInvalidWatchlistItem
	}

	action, err := s.actions.GetAction(ctx, actionID)
	if err != nil {
		return Item{}, err
	}
	if !action.Active {
		return Item{}, ErrWatchlistItemNotFound
	}

	return s.repo.Upsert(ctx, Item{
		UserID:    userID,
		ActionID:  action.ID,
		CreatedAt: time.Now().UTC(),
	})
}

func (s *Service) List(ctx context.Context, userID string) ([]Entry, error) {
	userID = strings.TrimSpace(userID)
	items, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	priceBySymbol := make(map[string]alerts.PriceSnapshot)
	if s.feed != nil {
		snapshots, err := s.feed.List(ctx)
		if err != nil {
			return nil, err
		}
		for _, snapshot := range snapshots {
			priceBySymbol[strings.ToUpper(strings.TrimSpace(snapshot.Symbol))] = snapshot
		}
	}

	byAction := make(map[string]int)
	if s.alerts != nil {
		userAlerts, err := s.alerts.ListByUser(ctx, userID)
		if err != nil {
			return nil, err
		}
		for _, alert := range userAlerts {
			if alert.Status == alerts.AlertStatusOpen {
				byAction[alert.ActionID]++
			}
		}
	}

	out := make([]Entry, 0, len(items))
	for _, item := range items {
		action, err := s.actions.GetAction(ctx, item.ActionID)
		if err != nil {
			return nil, err
		}

		entry := Entry{
			ActionID:        action.ID,
			Symbol:          action.Symbol,
			Name:            action.Name,
			Exchange:        action.Exchange,
			OpenAlertsCount: byAction[action.ID],
			CreatedAt:       item.CreatedAt,
		}

		if snapshot, ok := priceBySymbol[strings.ToUpper(strings.TrimSpace(action.Symbol))]; ok {
			price := snapshot.Price
			entry.CurrentPrice = &price
			if !snapshot.ObservedAt.IsZero() {
				at := snapshot.ObservedAt.UTC()
				entry.LastPriceAt = &at
			}
		}

		out = append(out, entry)
	}

	return out, nil
}

func (s *Service) Remove(ctx context.Context, userID, actionID string) error {
	userID = strings.TrimSpace(userID)
	actionID = strings.TrimSpace(actionID)
	if userID == "" || actionID == "" {
		return ErrInvalidWatchlistItem
	}

	if err := s.repo.Delete(ctx, userID, actionID); err != nil {
		return err
	}
	if s.alerts != nil {
		_, _ = s.alerts.DeleteByUserAndAction(ctx, userID, actionID)
	}
	return nil
}

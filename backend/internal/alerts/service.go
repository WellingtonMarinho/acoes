package alerts

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"ideacoes/backend/internal/actions"
)

type Notifier interface {
	Notify(ctx context.Context, alert Alert, marketPrice float64) error
}

type DeviceResolver interface {
	Resolve(ctx context.Context, userID string) (string, bool, error)
}

type SymbolRegistrar interface {
	RegisterSymbol(ctx context.Context, symbol string) error
}

type ActionResolver interface {
	GetAction(ctx context.Context, id string) (actions.Action, error)
}

type WatchlistRegistrar interface {
	Upsert(ctx context.Context, userID, actionID string) error
}

type Service struct {
	repo      Repository
	notifier  Notifier
	devices   DeviceResolver
	symbols   SymbolRegistrar
	watchlist WatchlistRegistrar
	actions   ActionResolver
}

func NewService(repo Repository, notifier Notifier, devices DeviceResolver) *Service {
	return &Service{repo: repo, notifier: notifier, devices: devices}
}

func NewServiceWithSymbolRegistrar(repo Repository, notifier Notifier, devices DeviceResolver, symbols SymbolRegistrar) *Service {
	return &Service{repo: repo, notifier: notifier, devices: devices, symbols: symbols}
}

func NewServiceWithActionResolver(repo Repository, notifier Notifier, devices DeviceResolver, symbols SymbolRegistrar, watchlist WatchlistRegistrar, actions ActionResolver) *Service {
	return &Service{repo: repo, notifier: notifier, devices: devices, symbols: symbols, watchlist: watchlist, actions: actions}
}

func (s *Service) CreateAlert(ctx context.Context, alert Alert) (Alert, error) {
	alert.ActionID = strings.TrimSpace(alert.ActionID)
	if alert.ActionID == "" || alert.TargetPrice <= 0 {
		return Alert{}, ErrInvalidAlert
	}
	if alert.Direction != DirectionAbove && alert.Direction != DirectionBelow {
		return Alert{}, ErrInvalidAlert
	}
	if s.actions != nil {
		action, err := s.actions.GetAction(ctx, alert.ActionID)
		if err != nil {
			return Alert{}, err
		}
		if !action.Active {
			return Alert{}, ErrInvalidAlert
		}
		alert.Symbol = strings.ToUpper(strings.TrimSpace(action.Symbol))
		alert.ActionName = action.Name
		if alert.Symbol == "" {
			return Alert{}, ErrInvalidAlert
		}
	}
	if s.watchlist != nil && strings.TrimSpace(alert.UserID) != "" {
		if err := s.watchlist.Upsert(ctx, alert.UserID, alert.ActionID); err != nil {
			return Alert{}, err
		}
	}
	if alert.DeviceToken == "" && s.devices != nil && strings.TrimSpace(alert.UserID) != "" {
		if token, ok, err := s.devices.Resolve(ctx, alert.UserID); err != nil {
			return Alert{}, err
		} else if ok {
			alert.DeviceToken = token
		}
	}
	if s.symbols != nil {
		if err := s.symbols.RegisterSymbol(ctx, alert.Symbol); err != nil {
			return Alert{}, err
		}
	}

	alert.ID = newID()
	alert.Status = AlertStatusOpen
	alert.CreatedAt = time.Now().UTC()
	alert.UpdatedAt = alert.CreatedAt
	return s.repo.Create(ctx, alert)
}

func (s *Service) ListAlerts(ctx context.Context) ([]Alert, error) {
	return s.repo.List(ctx)
}

func (s *Service) ListAlertsByUser(ctx context.Context, userID string) ([]Alert, error) {
	return s.repo.ListByUser(ctx, strings.TrimSpace(userID))
}

func (s *Service) UpdateAlert(ctx context.Context, userID, id string, update AlertUpdate) (Alert, error) {
	userID = strings.TrimSpace(userID)
	id = strings.TrimSpace(id)
	if userID == "" || id == "" || update.TargetPrice <= 0 || (update.Direction != DirectionAbove && update.Direction != DirectionBelow) {
		return Alert{}, ErrInvalidAlert
	}

	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return Alert{}, err
	}
	if existing.UserID != userID {
		return Alert{}, ErrAlertNotFound
	}
	if existing.Status != AlertStatusOpen {
		return Alert{}, ErrAlertNotEditable
	}

	existing.TargetPrice = update.TargetPrice
	existing.Direction = update.Direction
	existing.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, existing)
}

func (s *Service) DeleteAlert(ctx context.Context, userID, id string) error {
	userID = strings.TrimSpace(userID)
	id = strings.TrimSpace(id)
	if userID == "" || id == "" {
		return ErrInvalidAlert
	}

	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if existing.UserID != userID {
		return ErrAlertNotFound
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) CheckPrices(ctx context.Context, snapshots []PriceSnapshot) ([]Alert, error) {
	var triggered []Alert

	for _, snapshot := range snapshots {
		symbol := strings.ToUpper(strings.TrimSpace(snapshot.Symbol))
		if symbol == "" || snapshot.Price <= 0 {
			continue
		}

		alertsForSymbol, err := s.repo.ListOpenBySymbol(ctx, symbol)
		if err != nil {
			return nil, err
		}

		for _, alert := range alertsForSymbol {
			if !shouldTrigger(alert, snapshot.Price) {
				continue
			}

			triggeredAt := time.Now().UTC()
			updated, err := s.repo.MarkTriggered(ctx, alert.ID, triggeredAt)
			if err != nil {
				if errors.Is(err, ErrAlertNotFound) || errors.Is(err, ErrAlertNotEditable) {
					continue
				}
				return nil, err
			}

			if s.notifier != nil {
				if err := s.notifier.Notify(ctx, updated, snapshot.Price); err != nil {
					return nil, err
				}
			}
			triggered = append(triggered, updated)
		}
	}

	return triggered, nil
}

func shouldTrigger(alert Alert, marketPrice float64) bool {
	switch alert.Direction {
	case DirectionAbove:
		return marketPrice >= alert.TargetPrice
	case DirectionBelow:
		return marketPrice <= alert.TargetPrice
	default:
		return false
	}
}

type LogNotifier struct {
	logger interface{ Printf(string, ...any) }
}

func NewLogNotifier(logger interface{ Printf(string, ...any) }) *LogNotifier {
	return &LogNotifier{logger: logger}
}

func (n *LogNotifier) Notify(ctx context.Context, alert Alert, marketPrice float64) error {
	_ = ctx
	n.logger.Printf("alert triggered id=%s symbol=%s target=%.2f price=%.2f direction=%s", alert.ID, alert.Symbol, alert.TargetPrice, marketPrice, alert.Direction)
	return nil
}

func newID() string {
	return fmt.Sprintf("%d", time.Now().UTC().UnixNano())
}

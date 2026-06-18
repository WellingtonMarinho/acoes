package alerts

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrInvalidAlert = errors.New("invalid alert")

type Notifier interface {
	Notify(ctx context.Context, alert Alert, marketPrice float64) error
}

type DeviceResolver interface {
	Resolve(ctx context.Context, userID string) (string, bool, error)
}

type Service struct {
	repo     Repository
	notifier Notifier
	devices  DeviceResolver
}

func NewService(repo Repository, notifier Notifier, devices DeviceResolver) *Service {
	return &Service{repo: repo, notifier: notifier, devices: devices}
}

func (s *Service) CreateAlert(ctx context.Context, alert Alert) (Alert, error) {
	alert.Symbol = strings.ToUpper(strings.TrimSpace(alert.Symbol))
	if alert.Symbol == "" || alert.TargetPrice <= 0 {
		return Alert{}, ErrInvalidAlert
	}
	if alert.Direction != DirectionAbove && alert.Direction != DirectionBelow {
		return Alert{}, ErrInvalidAlert
	}
	if alert.DeviceToken == "" && s.devices != nil && strings.TrimSpace(alert.UserID) != "" {
		if token, ok, err := s.devices.Resolve(ctx, alert.UserID); err != nil {
			return Alert{}, err
		} else if ok {
			alert.DeviceToken = token
		}
	}

	alert.ID = newID()
	alert.Status = AlertStatusOpen
	alert.CreatedAt = time.Now().UTC()
	return s.repo.Create(ctx, alert)
}

func (s *Service) ListAlerts(ctx context.Context) ([]Alert, error) {
	return s.repo.List(ctx)
}

func (s *Service) CheckPrices(ctx context.Context, snapshots []PriceSnapshot) ([]Alert, error) {
	var triggered []Alert

	for _, snapshot := range snapshots {
		symbol := strings.ToUpper(strings.TrimSpace(snapshot.Symbol))
		if symbol == "" {
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
				return nil, err
			}

			if err := s.notifier.Notify(ctx, updated, snapshot.Price); err != nil {
				return nil, err
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

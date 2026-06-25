package pricefeed

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"ideacoes/backend/internal/alerts"
)

type TwelveDataFeed struct {
	client    *http.Client
	baseURL   string
	apiKey    string
	mu        sync.RWMutex
	symbols   map[string]struct{}
	overrides map[string]alerts.PriceSnapshot
}

func NewTwelveDataFeed(client *http.Client, baseURL, apiKey string) *TwelveDataFeed {
	if client == nil {
		client = http.DefaultClient
	}
	if baseURL == "" {
		baseURL = "https://api.twelvedata.com"
	}
	return &TwelveDataFeed{
		client:    client,
		baseURL:   strings.TrimRight(baseURL, "/"),
		apiKey:    apiKey,
		symbols:   make(map[string]struct{}),
		overrides: make(map[string]alerts.PriceSnapshot),
	}
}

func (f *TwelveDataFeed) List(ctx context.Context) ([]alerts.PriceSnapshot, error) {
	f.mu.RLock()
	symbols := make([]string, 0, len(f.symbols)+len(f.overrides))
	for symbol := range f.symbols {
		symbols = append(symbols, symbol)
	}
	for symbol := range f.overrides {
		if _, ok := f.symbols[symbol]; !ok {
			symbols = append(symbols, symbol)
		}
	}
	f.mu.RUnlock()

	if len(symbols) == 0 {
		return []alerts.PriceSnapshot{}, nil
	}

	out := make([]alerts.PriceSnapshot, 0, len(symbols))
	for _, symbol := range symbols {
		f.mu.RLock()
		if snapshot, ok := f.overrides[symbol]; ok {
			out = append(out, snapshot)
			f.mu.RUnlock()
			continue
		}
		f.mu.RUnlock()

		snapshot, err := f.fetchQuote(ctx, symbol)
		if err != nil {
			return nil, err
		}
		out = append(out, snapshot)
	}

	sort.Slice(out, func(i, j int) bool {
		return strings.ToUpper(out[i].Symbol) < strings.ToUpper(out[j].Symbol)
	})
	return out, nil
}

func (f *TwelveDataFeed) Upsert(ctx context.Context, snapshot alerts.PriceSnapshot) error {
	_ = ctx
	snapshot.Symbol = strings.ToUpper(strings.TrimSpace(snapshot.Symbol))
	if snapshot.Symbol == "" || snapshot.Price <= 0 {
		return fmt.Errorf("invalid price snapshot")
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	f.symbols[snapshot.Symbol] = struct{}{}
	f.overrides[snapshot.Symbol] = snapshot
	return nil
}

func (f *TwelveDataFeed) RegisterSymbol(ctx context.Context, symbol string) error {
	_ = ctx
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	f.symbols[symbol] = struct{}{}
	return nil
}

func (f *TwelveDataFeed) fetchQuote(ctx context.Context, symbol string) (alerts.PriceSnapshot, error) {
	endpoint, err := url.Parse(f.baseURL + "/time_series")
	if err != nil {
		return alerts.PriceSnapshot{}, err
	}

	query := endpoint.Query()
	query.Set("symbol", symbol)
	query.Set("interval", "1min")
	query.Set("outputsize", "1")
	query.Set("apikey", f.apiKey)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return alerts.PriceSnapshot{}, err
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return alerts.PriceSnapshot{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return alerts.PriceSnapshot{}, fmt.Errorf("twelvedata status %s", resp.Status)
	}

	var payload struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Meta    struct {
			Symbol string `json:"symbol"`
		} `json:"meta"`
		Values []struct {
			Close    string `json:"close"`
			Datetime string `json:"datetime"`
		} `json:"values"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return alerts.PriceSnapshot{}, err
	}
	if len(payload.Values) == 0 {
		return alerts.PriceSnapshot{}, fmt.Errorf("twelvedata returned no values for %s", symbol)
	}

	price, err := parseDecimal(payload.Values[0].Close)
	if err != nil {
		return alerts.PriceSnapshot{}, err
	}

	var observedAt time.Time
	if parsed, err := parseQuoteTime(payload.Values[0].Datetime); err == nil {
		observedAt = parsed
	}

	return alerts.PriceSnapshot{
		Symbol:     strings.ToUpper(strings.TrimSpace(symbol)),
		Price:      price,
		ObservedAt: observedAt,
	}, nil
}

func parseQuoteTime(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, fmt.Errorf("empty quote datetime")
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, raw); err == nil {
			return parsed.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported quote datetime %q", raw)
}

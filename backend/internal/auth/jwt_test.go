package auth

import (
	"context"
	"testing"
	"time"
)

func TestSignAndParse(t *testing.T) {
	now := time.Date(2026, time.June, 18, 12, 0, 0, 0, time.UTC)

	token, err := Sign("user-123", "secret", time.Hour, now)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	claims, err := Parse(token, "secret", now.Add(30*time.Minute))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if claims.UserID != "user-123" {
		t.Fatalf("expected user-123, got %q", claims.UserID)
	}
}

func TestParseRejectsExpiredToken(t *testing.T) {
	now := time.Date(2026, time.June, 18, 12, 0, 0, 0, time.UTC)

	token, err := Sign("user-123", "secret", time.Hour, now)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	if _, err := Parse(token, "secret", now.Add(2*time.Hour)); err == nil {
		t.Fatal("expected Parse() to fail for expired token")
	}
}

func TestBearerToken(t *testing.T) {
	token, err := BearerToken("Bearer abc.def.ghi")
	if err != nil {
		t.Fatalf("BearerToken() error = %v", err)
	}
	if token != "abc.def.ghi" {
		t.Fatalf("expected token abc.def.ghi, got %q", token)
	}
}

func TestContextClaimsRoundTrip(t *testing.T) {
	ctx := WithUserID(context.Background(), "user-123")
	userID, ok := UserIDFromContext(ctx)
	if !ok {
		t.Fatal("expected user id in context")
	}
	if userID != "user-123" {
		t.Fatalf("expected user-123, got %q", userID)
	}
}

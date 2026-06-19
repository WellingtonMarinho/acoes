package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrMissingToken = errors.New("missing token")
	ErrInvalidToken = errors.New("invalid token")
	ErrInvalidJWT   = errors.New("invalid jwt")
)

type userIDContextKey struct{}
type claimsContextKey struct{}

type Claims struct {
	UserID string `json:"user_id"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey{}, userID)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey{}).(string)
	return userID, ok && userID != ""
}

func Sign(userID, secret string, ttl time.Duration, now time.Time) (string, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return "", ErrInvalidToken
	}
	if strings.TrimSpace(secret) == "" {
		return "", errors.New("jwt secret is required")
	}

	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	claims := Claims{
		UserID: userID,
		Iat:    now.UTC().Unix(),
		Exp:    now.UTC().Add(ttl).Unix(),
	}

	headerPart, err := encodeJSON(header)
	if err != nil {
		return "", err
	}
	claimsPart, err := encodeJSON(claims)
	if err != nil {
		return "", err
	}

	unsigned := headerPart + "." + claimsPart
	sig := sign(unsigned, secret)
	return unsigned + "." + sig, nil
}

func Parse(token, secret string, now time.Time) (Claims, error) {
	if strings.TrimSpace(secret) == "" {
		return Claims{}, errors.New("jwt secret is required")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, ErrInvalidJWT
	}

	unsigned := parts[0] + "." + parts[1]
	expected := sign(unsigned, secret)
	if !hmac.Equal([]byte(parts[2]), []byte(expected)) {
		return Claims{}, ErrInvalidJWT
	}

	var header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}
	if err := decodeJSON(parts[0], &header); err != nil {
		return Claims{}, ErrInvalidJWT
	}
	if header.Alg != "HS256" || header.Typ != "JWT" {
		return Claims{}, ErrInvalidJWT
	}

	var claims Claims
	if err := decodeJSON(parts[1], &claims); err != nil {
		return Claims{}, ErrInvalidJWT
	}
	if claims.UserID == "" || claims.Exp == 0 {
		return Claims{}, ErrInvalidJWT
	}
	if now.UTC().Unix() > claims.Exp {
		return Claims{}, ErrInvalidJWT
	}

	return claims, nil
}

func encodeJSON(v any) (string, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func decodeJSON(part string, v any) error {
	raw, err := base64.RawURLEncoding.DecodeString(part)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, v)
}

func sign(unsigned, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(unsigned))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func BearerToken(header string) (string, error) {
	header = strings.TrimSpace(header)
	if header == "" {
		return "", ErrMissingToken
	}
	if !strings.HasPrefix(strings.ToLower(header), "bearer ") {
		return "", ErrInvalidToken
	}
	token := strings.TrimSpace(header[7:])
	if token == "" {
		return "", ErrMissingToken
	}
	return token, nil
}

func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey{}).(Claims)
	return claims, ok
}

func WithClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey{}, claims)
}

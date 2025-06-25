package token

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type service struct {
	secret []byte
	expiry time.Duration
}

func newService(secret string, expiryHours int) Service {
	return &service{
		secret: []byte(secret),
		expiry: time.Duration(expiryHours) * time.Hour,
	}
}

func (s *service) GenerateToken(ctx context.Context, userID uint, email string) (string, error) {
	claims := CustomClaims{
		Sub:   userID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *service) ParseToken(ctx context.Context, tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	return claims, nil
}

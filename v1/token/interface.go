package token

import "context"

// Service defines the interface for token handling.
type Service interface {
	// GenerateToken creates a JWT string with the given user ID and email
	GenerateToken(ctx context.Context, userID uint, email string) (string, error)

	// ParseToken validates the token and returns its claims if valid
	ParseToken(ctx context.Context, tokenStr string) (*CustomClaims, error)
}

func NewService(secret string, expiryHours int) Service {
	return newService(secret, expiryHours)
}

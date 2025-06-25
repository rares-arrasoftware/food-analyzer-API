package auth

import (
	"github.com/rares-arrasoftware/food-analyzer-api/v1/config"
	"github.com/rares-arrasoftware/food-analyzer-api/v1/token"
)

// Service defines the authentication service interface.
// It provides methods for registering and logging in users.
type Service interface {
	// Register creates a new user account and returns an auth token.
	Register(req RegisterRequest) (AuthResponse, error)

	// Login authenticates a user and returns an auth token.
	Login(req LoginRequest) (AuthResponse, error)

	// ValidateTokenFromHeader extracts the token from the Authorization header,
	// validates it, and returns the decoded claims if valid.
	// Returns an error if the token is missing, malformed, expired, or invalid.
	ValidateTokenFromHeader(authHeader string) (*token.CustomClaims, error)
}

func NewService(cfg config.Config) Service {
	return newService(cfg)
}

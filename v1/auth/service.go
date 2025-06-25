package auth

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/rares-arrasoftware/food-analyzer-api/v1/config"
	"github.com/rares-arrasoftware/food-analyzer-api/v1/database"
	"github.com/rares-arrasoftware/food-analyzer-api/v1/models"
	"github.com/rares-arrasoftware/food-analyzer-api/v1/token"
)

type ServiceImpl struct {
	users  database.Database[models.User]
	tokens token.Service
}

func newService(cfg config.Config) Service {
	userDB, err := database.NewDatabase(cfg.DatabaseDSN, models.User{})
	if err != nil {
		log.Fatalf("failed to initialize user database: %v", err)
	}

	return &ServiceImpl{
		users:  userDB,
		tokens: token.NewService(cfg.JWTSecret, cfg.JWTExpiry),
	}
}

func (s *ServiceImpl) Register(req RegisterRequest) (resp AuthResponse, err error) {
	if err = validateCredentialsInput(req.Email, req.Password); err != nil {
		return
	}

	if exists, _ := s.users.GetByField("email", req.Email); exists != nil {
		err = errors.New("email already exists")
		return
	}

	user, err := newUserFromRequest(req)
	if err != nil {
		return
	}

	if err = s.users.Create(user); err != nil {
		return
	}

	return s.generateAuthResponse(user)
}

func (s *ServiceImpl) Login(req LoginRequest) (resp AuthResponse, err error) {
	if err = validateCredentialsInput(req.Email, req.Password); err != nil {
		return
	}

	user, err := s.users.GetByField("email", req.Email)
	if err != nil {
		err = errors.New("invalid email")
		return
	}

	if err = ComparePassword(user.Password, req.Password); err != nil {
		err = errors.New("invalid password")
		return
	}

	return s.generateAuthResponse(*user)
}

func (s *ServiceImpl) ValidateTokenFromHeader(authHeader string) (*token.CustomClaims, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("invalid authorization header")
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	return s.tokens.ParseToken(context.Background(), tokenStr)
}

// Helper functions

func validateCredentialsInput(email, password string) error {
	// Here we can add more rules about password/email like password length etc.
	// also we can use an external lib(e.g: go-playground/validator/v10) if we want complex validation.

	if email == "" || password == "" {
		return errors.New("email and password are required")
	}

	return nil
}

func newUserFromRequest(req RegisterRequest) (models.User, error) {
	hashed, err := HashPassword(req.Password)
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		Email:     req.Email,
		Password:  hashed,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}, nil
}

func (s *ServiceImpl) generateAuthResponse(user models.User) (resp AuthResponse, err error) {
	tokenStr, err := s.tokens.GenerateToken(context.Background(), user.ID, user.Email)
	if err != nil {
		return
	}
	return AuthResponse{Token: tokenStr}, nil
}

package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/rares-arrasoftware/food-analyzer-api/v1/models"
	"github.com/rares-arrasoftware/food-analyzer-api/v1/token"
)

// --- Mock DB ---
type mockUserDB struct {
	users map[any]models.User
}

func (m *mockUserDB) GetByField(field string, value any) (*models.User, error) {
	u, ok := m.users[value]
	if !ok {
		return nil, errors.New("not found")
	}
	return &u, nil
}

func (m *mockUserDB) Create(user models.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("already exists")
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserDB) Delete(uint) error {
	return nil
}

func (m *mockUserDB) GetByID(uint) (*models.User, error) {
	return nil, nil
}

func (m *mockUserDB) Update(user models.User) error {
	return nil
}

// --- Mock Token Service ---
type mockTokenService struct{}

func (m *mockTokenService) GenerateToken(ctx context.Context, userID uint, email string) (string, error) {
	return "mock.token.string", nil
}

func (m *mockTokenService) ParseToken(ctx context.Context, tokenStr string) (*token.CustomClaims, error) {
	if tokenStr == "invalid" {
		return nil, errors.New("invalid token")
	}
	return &token.CustomClaims{Email: "test@example.com"}, nil
}

// --- Test Constructor ---
func newMockedService(userDB *mockUserDB, tokenSvc *mockTokenService) *ServiceImpl {
	return &ServiceImpl{
		users:  userDB,
		tokens: tokenSvc,
	}
}

// --- Tests ---
func TestRegister_Success(t *testing.T) {
	svc := newMockedService(&mockUserDB{users: map[any]models.User{}}, &mockTokenService{})

	req := RegisterRequest{
		Email:     "john@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	resp, err := svc.Register(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected a token, got empty string")
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc := newMockedService(&mockUserDB{users: map[any]models.User{
		"john@example.com": {Email: "john@example.com"},
	}}, &mockTokenService{})

	req := RegisterRequest{
		Email:     "john@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	_, err := svc.Register(req)
	if err == nil {
		t.Fatal("expected error for duplicate email, got nil")
	}
}

func TestLogin_Success(t *testing.T) {
	hashed, _ := HashPassword("password123")
	user := models.User{
		Email:    "john@example.com",
		Password: hashed,
	}

	svc := newMockedService(&mockUserDB{users: map[any]models.User{"john@example.com": user}}, &mockTokenService{})

	req := LoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	}

	resp, err := svc.Login(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected token, got empty string")
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	hashed, _ := HashPassword("password123")
	user := models.User{
		Email:    "john@example.com",
		Password: hashed,
	}

	svc := newMockedService(&mockUserDB{users: map[any]models.User{"john@example.com": user}}, &mockTokenService{})

	req := LoginRequest{
		Email:    "john@example.com",
		Password: "wrongpass",
	}

	_, err := svc.Login(req)
	if err == nil {
		t.Fatal("expected error for invalid password, got nil")
	}
}

func TestValidateTokenFromHeader(t *testing.T) {
	svc := newMockedService(nil, &mockTokenService{})

	claims, err := svc.ValidateTokenFromHeader("Bearer mock.token.string")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if claims.Email != "test@example.com" {
		t.Errorf("expected email to be 'test@example.com', got %s", claims.Email)
	}
}

func TestValidateTokenFromHeader_InvalidFormat(t *testing.T) {
	svc := newMockedService(nil, &mockTokenService{})

	_, err := svc.ValidateTokenFromHeader("InvalidHeader")
	if err == nil {
		t.Fatal("expected error for invalid header format, got nil")
	}
}

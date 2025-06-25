package token

import (
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Sub   uint   `json:"sub"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

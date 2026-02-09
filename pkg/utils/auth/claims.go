package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the custom JWT claims embedded in every access token.
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Roles  []string  `json:"roles"`
	jwt.RegisteredClaims
}

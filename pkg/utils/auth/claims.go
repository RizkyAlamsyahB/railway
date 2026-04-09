package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the custom JWT claims embedded in every access token.
type Claims struct {
	UserID            uuid.UUID  `json:"user_id"`
	Email             string     `json:"email"`
	Role              string     `json:"role"`
	VendorID          *uuid.UUID `json:"vendor_id,omitempty"`
	PasswordChangedAt *time.Time `json:"pwd_changed_at,omitempty"`
	jwt.RegisteredClaims
}

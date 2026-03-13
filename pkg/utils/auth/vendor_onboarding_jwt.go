package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const VendorOnboardingTokenType = "vendor_onboarding"

// VendorOnboardingClaims represents the JWT claims used during vendor onboarding.
type VendorOnboardingClaims struct {
	OnboardingID uuid.UUID `json:"onboarding_id"`
	Email        string    `json:"email"`
	Type         string    `json:"typ"`
	jwt.RegisteredClaims
}

// GenerateVendorOnboardingToken creates a signed JWT for vendor onboarding steps.
func GenerateVendorOnboardingToken(onboardingID uuid.UUID, email string, secret string, expiry time.Duration, issuer string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(expiry)

	claims := VendorOnboardingClaims{
		OnboardingID: onboardingID,
		Email:        email,
		Type:         VendorOnboardingTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   onboardingID.String(),
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign vendor onboarding token: %w", err)
	}

	return signed, expiresAt, nil
}

// ValidateVendorOnboardingToken parses and validates a vendor onboarding JWT string.
func ValidateVendorOnboardingToken(tokenString string, secret string) (*VendorOnboardingClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &VendorOnboardingClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*VendorOnboardingClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	if claims.Type != VendorOnboardingTokenType {
		return nil, fmt.Errorf("invalid token type")
	}

	return claims, nil
}

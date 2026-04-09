package domain

import "context"

// GoogleUserProfile holds verified claims extracted from a Google ID token.
type GoogleUserProfile struct {
	Subject       string
	Email         string
	Name          string
	PictureURL    string
	EmailVerified bool
}

// GoogleOAuthVerifier verifies Google ID tokens.
type GoogleOAuthVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*GoogleUserProfile, error)
}

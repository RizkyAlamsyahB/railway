package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/config"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type googleTokenInfoResponse struct {
	Audience      string `json:"aud"`
	Issuer        string `json:"iss"`
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type googleVerifier struct {
	client       *http.Client
	clientID     string
	tokenInfoURL string
}

// NewGoogleVerifier creates a Google ID token verifier.
func NewGoogleVerifier(cfg config.GoogleOAuthConfig) domain.GoogleOAuthVerifier {
	timeout := time.Duration(cfg.TimeoutSec) * time.Second
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &googleVerifier{
		client:       &http.Client{Timeout: timeout},
		clientID:     strings.TrimSpace(cfg.ClientID),
		tokenInfoURL: strings.TrimSpace(cfg.TokenInfoURL),
	}
}

func (v *googleVerifier) VerifyIDToken(ctx context.Context, idToken string) (*domain.GoogleUserProfile, error) {
	if strings.TrimSpace(v.clientID) == "" {
		return nil, fmt.Errorf("google oauth client id is not configured")
	}
	idToken = strings.Join(strings.Fields(idToken), "")
	if idToken == "" {
		return nil, fmt.Errorf("id token is required")
	}

	endpoint := strings.TrimRight(v.tokenInfoURL, "?") + "?id_token=" + url.QueryEscape(idToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := v.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call google tokeninfo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("google tokeninfo rejected token: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var payload googleTokenInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode google tokeninfo response: %w", err)
	}

	if payload.Audience != v.clientID {
		return nil, fmt.Errorf("google token audience mismatch")
	}
	if payload.Issuer != "accounts.google.com" && payload.Issuer != "https://accounts.google.com" {
		return nil, fmt.Errorf("google token issuer is invalid")
	}

	return &domain.GoogleUserProfile{
		Subject:       payload.Subject,
		Email:         strings.TrimSpace(payload.Email),
		Name:          strings.TrimSpace(payload.Name),
		PictureURL:    strings.TrimSpace(payload.Picture),
		EmailVerified: strings.EqualFold(payload.EmailVerified, "true"),
	}, nil
}

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase"
)

type stubOTPUseCase struct {
	requestOTP      func(context.Context, domain.RequestOTPInput) (*domain.RequestOTPResult, error)
	verifyOTP       func(context.Context, domain.VerifyOTPInput) (*domain.VerifyOTPResult, error)
	verifyProofFunc func(context.Context, string, string) (*domain.OTPProofClaims, error)
}

func (s stubOTPUseCase) RequestOTP(ctx context.Context, input domain.RequestOTPInput) (*domain.RequestOTPResult, error) {
	return s.requestOTP(ctx, input)
}

func (s stubOTPUseCase) VerifyOTP(ctx context.Context, input domain.VerifyOTPInput) (*domain.VerifyOTPResult, error) {
	return s.verifyOTP(ctx, input)
}

func (s stubOTPUseCase) VerifyProofToken(ctx context.Context, token string, expectedPurpose string) (*domain.OTPProofClaims, error) {
	if s.verifyProofFunc == nil {
		return nil, nil
	}
	return s.verifyProofFunc(ctx, token, expectedPurpose)
}

func TestOTPHandlerRequestOTP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		h := NewOTPHandler(stubOTPUseCase{
			requestOTP: func(_ context.Context, input domain.RequestOTPInput) (*domain.RequestOTPResult, error) {
				if input.Email != "user@example.com" {
					t.Fatalf("unexpected email %s", input.Email)
				}
				return &domain.RequestOTPResult{
					ExpiresAt:     time.Now().Add(5 * time.Minute),
					CooldownUntil: time.Now().Add(1 * time.Minute),
				}, nil
			},
			verifyOTP: func(context.Context, domain.VerifyOTPInput) (*domain.VerifyOTPResult, error) { return nil, nil },
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/otp/request", strings.NewReader(`{"email":"user@example.com","purpose":"email_verification"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		h.RequestOTP(c)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", w.Code)
		}

		var env handlerEnvelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if env.Message != "OTP generated successfully" {
			t.Fatalf("unexpected message: %s", env.Message)
		}
	})

	t.Run("mapped usecase error", func(t *testing.T) {
		h := NewOTPHandler(stubOTPUseCase{
			requestOTP: func(context.Context, domain.RequestOTPInput) (*domain.RequestOTPResult, error) {
				return nil, usecase.ErrOTPResendTooSoon
			},
			verifyOTP: func(context.Context, domain.VerifyOTPInput) (*domain.VerifyOTPResult, error) { return nil, nil },
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/otp/request", strings.NewReader(`{"email":"user@example.com","purpose":"email_verification"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		h.RequestOTP(c)

		if w.Code != http.StatusTooManyRequests {
			t.Fatalf("expected status 429, got %d", w.Code)
		}
	})
}

func TestOTPHandlerVerifyOTP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		h := NewOTPHandler(stubOTPUseCase{
			requestOTP: func(context.Context, domain.RequestOTPInput) (*domain.RequestOTPResult, error) { return nil, nil },
			verifyOTP: func(_ context.Context, input domain.VerifyOTPInput) (*domain.VerifyOTPResult, error) {
				if input.Code != "123456" {
					t.Fatalf("unexpected code %s", input.Code)
				}
				return &domain.VerifyOTPResult{
					ProofToken:     "proof-token",
					ProofExpiresAt: time.Now().Add(10 * time.Minute),
					Email:          input.Email,
					Purpose:        input.Purpose,
				}, nil
			},
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/otp/verify", strings.NewReader(`{"email":"user@example.com","purpose":"forgot_password","code":"123456"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		h.VerifyOTP(c)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		var env handlerEnvelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if env.Message != "OTP verified successfully" {
			t.Fatalf("unexpected message: %s", env.Message)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		h := NewOTPHandler(stubOTPUseCase{
			requestOTP: func(context.Context, domain.RequestOTPInput) (*domain.RequestOTPResult, error) { return nil, nil },
			verifyOTP:  func(context.Context, domain.VerifyOTPInput) (*domain.VerifyOTPResult, error) { return nil, nil },
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/otp/verify", strings.NewReader(`{"email":`))
		c.Request.Header.Set("Content-Type", "application/json")

		h.VerifyOTP(c)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}
	})
}

package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
)

// Context keys for storing authenticated user claims.
const (
	ContextKeyClaims   = "auth_claims"
	ContextKeyUserID   = "auth_user_id"
	ContextKeyEmail    = "auth_email"
	ContextKeyRole     = "auth_role"
	ContextKeyVendorID = "auth_vendor_id"

	ContextKeyVendorOnboardingClaims = "vendor_onboarding_claims"
	ContextKeyVendorOnboardingID     = "vendor_onboarding_id"
)

// Auth returns a middleware that validates the JWT Bearer token from the
// Authorization header. On success, the decoded claims are stored in the
// Gin context for downstream handlers.
func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Abort(c, http.StatusUnauthorized, "missing authorization header", nil)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Abort(c, http.StatusUnauthorized, "invalid authorization format", nil)
			return
		}

		claims, err := auth.ValidateToken(parts[1], jwtSecret)
		if err != nil {
			response.Abort(c, http.StatusUnauthorized, "invalid or expired token", nil)
			return
		}

		c.Set(ContextKeyClaims, claims)
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyEmail, claims.Email)
		c.Set(ContextKeyRole, claims.Role)
		if claims.VendorID != nil {
			c.Set(ContextKeyVendorID, *claims.VendorID)
		}

		c.Next()
	}
}

// AuthVendorOnboarding validates a vendor onboarding JWT and stores onboarding claims in context.
func AuthVendorOnboarding(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Abort(c, http.StatusUnauthorized, "missing authorization header", nil)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Abort(c, http.StatusUnauthorized, "invalid authorization format", nil)
			return
		}

		claims, err := auth.ValidateVendorOnboardingToken(parts[1], jwtSecret)
		if err != nil {
			message := "invalid vendor onboarding token"
			if errors.Is(err, jwt.ErrTokenExpired) {
				message = "vendor onboarding token has expired"
			}
			response.Abort(c, http.StatusUnauthorized, message, nil)
			return
		}

		c.Set(ContextKeyVendorOnboardingClaims, claims)
		c.Set(ContextKeyVendorOnboardingID, claims.OnboardingID)
		c.Set(ContextKeyEmail, claims.Email)

		c.Next()
	}
}

// RequireRoles returns a middleware that checks whether the authenticated user
// has one of the specified roles. It must be used after the Auth
// middleware in the handler chain.
func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c *gin.Context) {
		claimsVal, exists := c.Get(ContextKeyClaims)
		if !exists {
			response.Abort(c, http.StatusUnauthorized, "authentication required", nil)
			return
		}

		claims, ok := claimsVal.(*auth.Claims)
		if !ok {
			response.Abort(c, http.StatusInternalServerError, "invalid claims in context", nil)
			return
		}

		if _, ok := allowed[claims.Role]; ok {
			c.Next()
			return
		}

		response.Abort(c, http.StatusForbidden, "insufficient permissions", nil)
	}
}

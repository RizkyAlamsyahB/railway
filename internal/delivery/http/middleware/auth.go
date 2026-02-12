package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
)

// Context keys for storing authenticated user claims.
const (
	ContextKeyClaims   = "auth_claims"
	ContextKeyUserID   = "auth_user_id"
	ContextKeyEmail    = "auth_email"
	ContextKeyRoles    = "auth_roles"
	ContextKeyVendorID = "auth_vendor_id"
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
		c.Set(ContextKeyRoles, claims.Roles)
		if claims.VendorID != nil {
			c.Set(ContextKeyVendorID, *claims.VendorID)
		}

		c.Next()
	}
}

// RequireRoles returns a middleware that checks whether the authenticated user
// has at least one of the specified roles. It must be used after the Auth
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

		for _, role := range claims.Roles {
			if _, ok := allowed[role]; ok {
				c.Next()
				return
			}
		}

		response.Abort(c, http.StatusForbidden, "insufficient permissions", nil)
	}
}

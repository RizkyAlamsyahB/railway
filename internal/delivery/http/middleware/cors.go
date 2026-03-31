package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const defaultCORSAllowHeaders = "Origin, Content-Type, Accept, Authorization, X-Request-ID"

func mergeAllowedHeaders(requestedHeaders string) string {
	seen := make(map[string]struct{})
	headers := make([]string, 0)

	appendHeaders := func(raw string) {
		for _, header := range strings.Split(raw, ",") {
			header = strings.TrimSpace(header)
			if header == "" {
				continue
			}

			key := strings.ToLower(header)
			if _, exists := seen[key]; exists {
				continue
			}

			seen[key] = struct{}{}
			headers = append(headers, header)
		}
	}

	appendHeaders(defaultCORSAllowHeaders)
	appendHeaders(requestedHeaders)

	return strings.Join(headers, ", ")
}

// CORS returns a middleware that sets Cross-Origin Resource Sharing headers.
// allowedOrigins is a comma-separated list of allowed origins (e.g. "https://example.com,https://app.example.com").
// Use "*" to allow all origins.
func CORS(allowedOrigins string) gin.HandlerFunc {
	allowed := parseAllowedOrigins(allowedOrigins)

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = "*"
		}

		// Check if origin is allowed.
		resOrigin := resolveOrigin(origin, allowed)

		c.Header("Access-Control-Allow-Origin", resOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", mergeAllowedHeaders(c.GetHeader("Access-Control-Request-Headers")))
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400")
		c.Header("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// parseAllowedOrigins splits a comma-separated origin string into a set.
func parseAllowedOrigins(raw string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, o := range strings.Split(raw, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			result[o] = struct{}{}
		}
	}
	return result
}

// resolveOrigin returns the origin to set in Access-Control-Allow-Origin.
func resolveOrigin(reqOrigin string, allowed map[string]struct{}) string {
	if _, ok := allowed["*"]; ok {
		return reqOrigin
	}
	if _, ok := allowed[reqOrigin]; ok {
		return reqOrigin
	}
	return ""
}

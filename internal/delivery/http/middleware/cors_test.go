package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSOptionsAllowsRequestedHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(CORS("*"))
	r.OPTIONS("/api/v1/users/login", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/users/login", nil)
	req.Header.Set("Origin", "http://127.0.0.1:5500")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	req.Header.Set("Access-Control-Request-Headers", "content-type,ngrok-skip-browser-warning,authorization")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://127.0.0.1:5500" {
		t.Fatalf("expected allow origin to echo request origin, got %q", got)
	}

	allowHeaders := strings.ToLower(w.Header().Get("Access-Control-Allow-Headers"))
	for _, header := range []string{"content-type", "authorization", "x-request-id", "ngrok-skip-browser-warning"} {
		if !strings.Contains(allowHeaders, header) {
			t.Fatalf("expected allow headers %q to contain %q", allowHeaders, header)
		}
	}
}

func TestMergeAllowedHeadersDeduplicatesValues(t *testing.T) {
	got := mergeAllowedHeaders("content-type, X-Request-ID, ngrok-skip-browser-warning, authorization")
	lower := strings.ToLower(got)

	if strings.Count(lower, "content-type") != 1 {
		t.Fatalf("expected content-type once, got %q", got)
	}

	if !strings.Contains(lower, "ngrok-skip-browser-warning") {
		t.Fatalf("expected ngrok header in %q", got)
	}
}

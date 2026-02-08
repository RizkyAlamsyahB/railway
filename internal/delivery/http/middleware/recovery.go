package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// Recovery returns a middleware that recovers from panics and returns a
// JSON envelope error response instead of crashing the server.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				response.Abort(c, http.StatusInternalServerError, "internal server error", nil)
			}
		}()
		c.Next()
	}
}

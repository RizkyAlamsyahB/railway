package websocket

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	gorillaws "github.com/gorilla/websocket"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
)

var upgrader = gorillaws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins — CORS is handled at the Gin level.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Handler handles GET /api/v1/ws
// The client must pass the JWT as a query param: ?token=<jwt>
// (WebSocket browser API does not support custom headers.)
type Handler struct {
	hub       *Hub
	jwtSecret string
}

// NewHandler creates a new WebSocket HTTP handler.
func NewHandler(hub *Hub, jwtSecret string) *Handler {
	return &Handler{hub: hub, jwtSecret: jwtSecret}
}

// Connect upgrades the HTTP connection to WebSocket and blocks until
// the client disconnects.
//
// GET /api/v1/ws?token=<jwt>
func (h *Handler) Connect(c *gin.Context) {
	// --- Auth ---
	tokenStr := strings.TrimSpace(c.Query("token"))
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "missing token"})
		return
	}

	claims, err := auth.ValidateToken(tokenStr, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid token"})
		return
	}

	// --- Upgrade ---
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// Upgrade already writes the error response.
		log.Printf("[ws] upgrade error for user %s: %v", claims.UserID, err)
		return
	}

	// ServeClient blocks until the connection closes.
	h.hub.ServeClient(conn, claims.UserID)
}

package realtime

import (
	"log/slog"
	"net/http"

	"freetokenspoker/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Handler upgrades HTTP connections to WebSockets and binds them to a user.
type Handler struct {
	hub      *Hub
	jwt      *auth.JWTManager
	log      *slog.Logger
	upgrader websocket.Upgrader
}

// NewHandler builds the WebSocket handler. allowedOrigins gates the upgrade.
func NewHandler(hub *Hub, jwt *auth.JWTManager, log *slog.Logger, allowedOrigins []string) *Handler {
	allow := map[string]struct{}{}
	for _, o := range allowedOrigins {
		allow[o] = struct{}{}
	}
	return &Handler{
		hub: hub,
		jwt: jwt,
		log: log,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if origin == "" {
					return true // non-browser clients (smoke tests)
				}
				_, ok := allow[origin]
				return ok
			},
		},
	}
}

// Upgrade authenticates via a token query param, upgrades the connection, and
// starts the read/write pumps.
func (h *Handler) Upgrade(c *gin.Context) {
	token := c.Query("token")
	claims, err := h.jwt.Verify(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid token", "errorCode": "UNAUTHORIZED"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Warn("websocket upgrade failed", "error", err)
		return
	}

	client := &Client{
		hub:    h.hub,
		conn:   conn,
		log:    h.log,
		send:   make(chan Envelope, sendBuffer),
		userID: claims.UserID,
		name:   claims.Name,
		email:  claims.Email,
		rooms:  map[string]struct{}{},
	}

	go client.writePump()
	go client.readPump()
}

package api

import (
	"freetokenspoker/internal/auth"
	"freetokenspoker/internal/config"
	"freetokenspoker/internal/handlers"
	"freetokenspoker/internal/middleware"
	"freetokenspoker/internal/realtime"
	"freetokenspoker/internal/repositories"
	"freetokenspoker/internal/services"

	"log/slog"

	"github.com/gin-gonic/gin"
)

// Deps bundles everything the router needs.
type Deps struct {
	Cfg   *config.Config
	Log   *slog.Logger
	Repos *repositories.Repositories
	JWT   *auth.JWTManager
	Hub   *realtime.Hub
}

// NewRouter builds the Gin engine with all middleware and routes.
func NewRouter(d Deps) *gin.Engine {
	if d.Cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Recovery(d.Log))
	r.Use(middleware.Logger(d.Log))
	r.Use(middleware.CORS(d.Cfg.CORSOrigins))

	limiter := middleware.NewRateLimiter(d.Cfg.RateLimitRPS, d.Cfg.RateLimitBurst)
	r.Use(limiter.Middleware())

	// Services.
	authSvc := services.NewAuthService(d.Repos.Users, d.JWT)
	roomSvc := services.NewRoomService(d.Repos.Rooms, d.Hub)
	taskSvc := services.NewTaskService(d.Repos.Tasks, d.Repos.Rooms, d.Repos.Votes, d.Repos.Finals, d.Hub)
	voteSvc := services.NewVoteService(d.Repos.Votes, d.Repos.Tasks, d.Repos.Rooms, d.Hub)
	histSvc := services.NewHistoryService(d.Repos.Rooms, d.Repos.Tasks, d.Repos.Votes, d.Repos.Finals)

	// Handlers.
	authH := handlers.NewAuthHandler(authSvc)
	roomH := handlers.NewRoomHandler(roomSvc)
	taskH := handlers.NewTaskHandler(taskSvc)
	voteH := handlers.NewVoteHandler(voteSvc)
	histH := handlers.NewHistoryHandler(histSvc)
	metaH := handlers.NewMetaHandler()

	// WebSocket upgrade endpoint (auth via token query param).
	wsHandler := realtime.NewHandler(d.Hub, d.JWT, d.Log, d.Cfg.CORSOrigins)
	r.GET("/ws", wsHandler.Upgrade)

	// Public endpoints.
	r.GET("/health", metaH.Health)

	api := r.Group("/api")
	{
		api.GET("/modes", metaH.Modes)
		api.POST("/auth/login", authH.Login)
		// Public invite preview: an unauthenticated invitee can see the room
		// name before deciding to sign in and join.
		api.GET("/invite/:code", roomH.Preview)

		// Authenticated endpoints.
		authed := api.Group("")
		authed.Use(middleware.Auth(d.JWT))
		{
			authed.GET("/me", authH.Me)

			authed.POST("/rooms", roomH.Create)
			authed.POST("/rooms/join", roomH.Join)
			authed.GET("/rooms/:id", roomH.Get)
			authed.GET("/rooms/:id/tasks", taskH.ListByRoom)
			authed.DELETE("/rooms/:id", roomH.Delete)

			authed.POST("/tasks", taskH.Create)
			authed.GET("/tasks/:id", taskH.Get)
			authed.PATCH("/tasks/:id/reveal", taskH.Reveal)
			authed.PATCH("/tasks/:id/final", taskH.Final)

			authed.POST("/votes", voteH.Submit)
			authed.PATCH("/votes", voteH.Submit)

			authed.GET("/history", histH.List)
		}
	}

	return r
}

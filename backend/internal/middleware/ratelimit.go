package middleware

import (
	"net/http"
	"sync"
	"time"

	"freetokenspoker/internal/dto"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter applies a per-client-IP token bucket. Lightweight and in-memory,
// sufficient for the MVP; swap for Redis when scaling horizontally.
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rps      rate.Limit
	burst    int
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewRateLimiter builds a limiter and starts a background cleanup loop.
func NewRateLimiter(rps, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rps:      rate.Limit(rps),
		burst:    burst,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) get(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	v, ok := rl.visitors[ip]
	if !ok {
		l := rate.NewLimiter(rl.rps, rl.burst)
		rl.visitors[ip] = &visitor{limiter: l, lastSeen: time.Now()}
		return l
	}
	v.lastSeen = time.Now()
	return v.limiter
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Middleware enforces the limit, returning 429 when exceeded.
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.get(c.ClientIP()).Allow() {
			dto.Fail(c, http.StatusTooManyRequests, "RATE_LIMITED", "too many requests, slow down")
			return
		}
		c.Next()
	}
}

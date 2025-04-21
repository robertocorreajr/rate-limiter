package handler

import (
	"net/http"
	"os"
	"strconv"

	"github.com/robertocorreajr/rate-limiter/internal/limiter"
)

type RateLimiterMiddleware struct {
	limiter *limiter.RateLimiterService
}

func NewRateLimiterMiddleware(l *limiter.RateLimiterService) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{limiter: l}
}

func (m *RateLimiterMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit, _ := strconv.Atoi(os.Getenv("DEFAULT_LIMIT"))
		if limit == 0 {
			limit = 10
		}
		blockTime, _ := strconv.Atoi(os.Getenv("BLOCK_TIME"))
		if blockTime == 0 {
			blockTime = 300
		}

		key := limiter.ExtractKey(r)
		allowed, err := m.limiter.Allow(key, limit, blockTime)
		if err != nil || !allowed {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

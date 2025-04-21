package limiter

import (
	"net/http"
	"strings"

	"github.com/robertocorreajr/rate-limiter/internal/storage"
)

type RateLimiterService struct {
	store storage.Storage
}

func NewRateLimiterService(s storage.Storage) *RateLimiterService {
	return &RateLimiterService{store: s}
}

func (r *RateLimiterService) Allow(key string, defaultLimit int, blockTime int) (bool, error) {
	if blocked, _ := r.store.IsBlocked(key); blocked {
		return false, nil
	}

	count, err := r.store.Increment(key)
	if err != nil {
		return false, err
	}

	ttl, _ := r.store.TTL(key)
	if ttl < 0 {
		_ = r.store.SetTTL(key, 1)
	}

	limit, err := r.store.GetLimit(key)
	if err != nil || limit == 0 {
		limit = defaultLimit
	}

	if count > limit {
		_ = r.store.Block(key, blockTime)
		return false, nil
	}
	return true, nil
}

func ExtractKey(r *http.Request) string {
	token := r.Header.Get("API_KEY")
	if token != "" {
		return "token:" + token
	}
	ip := strings.Split(r.RemoteAddr, ":")[0]
	return "ip:" + ip
}

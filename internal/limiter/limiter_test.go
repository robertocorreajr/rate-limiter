package limiter_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/robertocorreajr/rate-limiter/internal/handler"
	"github.com/robertocorreajr/rate-limiter/internal/limiter"
	"github.com/robertocorreajr/rate-limiter/internal/storage"
)

func TestRateLimiter_IPLimitExceeded(t *testing.T) {
	os.Setenv("DEFAULT_LIMIT", "3")
	os.Setenv("BLOCK_TIME", "1")

	store := storage.NewRedisStorage("localhost:6379")
	limiterService := limiter.NewRateLimiterService(store)
	middleware := handler.NewRateLimiterMiddleware(limiterService)

	h := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	for i := 1; i <= 4; i++ {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		if i <= 3 && rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", rr.Code)
		} else if i == 4 && rr.Code != http.StatusTooManyRequests {
			t.Errorf("Expected 429, got %d", rr.Code)
		}
	}

	time.Sleep(2 * time.Second) // Espera desbloqueio
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 after block expires, got %d", rr.Code)
	}
}

func TestRateLimiter_TokenOverridesIPLimit(t *testing.T) {
	os.Setenv("DEFAULT_LIMIT", "1")
	os.Setenv("BLOCK_TIME", "1")

	store := storage.NewRedisStorage("localhost:6379")
	store.SetLimit("token:abc123", 5)
	limiterService := limiter.NewRateLimiterService(store)
	middleware := handler.NewRateLimiterMiddleware(limiterService)

	h := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:9999"
	req.Header.Set("API_KEY", "abc123")

	for i := 1; i <= 5; i++ {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d at request %d", rr.Code, i)
		}
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Expected 429 after limit exceeded, got %d", rr.Code)
	}
}

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/robertocorreajr/rate-limiter/internal/handler"
	"github.com/robertocorreajr/rate-limiter/internal/limiter"
	"github.com/robertocorreajr/rate-limiter/internal/storage"
)

func main() {
	_ = godotenv.Load()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	db := storage.NewRedisStorage(redisAddr)
	limiterService := limiter.NewRateLimiterService(db)
	middleware := handler.NewRateLimiterMiddleware(limiterService)

	r := mux.NewRouter()
	r.Use(middleware.Middleware)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Request accepted")
	})

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}

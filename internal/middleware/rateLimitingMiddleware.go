package middleware

import (
	"net/http"

	"github.com/zhxauda9/StayMate/internal/myLogger"
	"golang.org/x/time/rate"
)

// Middleware for Rate Limitting
func RateLimiterMiddleware(next http.Handler, limiter *rate.Limiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Returns function that add middleware with Rate Limit
func RateLimiterMiddlewareFunc(limiter *rate.Limiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				myLogger.Log.Debug().Str("IP", r.RemoteAddr).Int("Status Code", http.StatusTooManyRequests).Msg("Rate limit exceeded")
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RateLimiterMiddlewareFun(limiter *rate.Limiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return RateLimiterMiddleware(next, limiter)
	}
}

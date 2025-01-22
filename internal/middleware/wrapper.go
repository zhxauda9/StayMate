package middleware

import "net/http"

func MiddlewareWrapper(wrapper func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return wrapper(next)
	}
}

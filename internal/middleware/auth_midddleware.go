package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func AuthAdminMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Authorization")
		if err != nil {
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			http.Error(w, "Access denied: Admins only", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func AuthUserMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the token from the cookie
		cookie, err := r.Cookie("Authorization")
		if err != nil {
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		_, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

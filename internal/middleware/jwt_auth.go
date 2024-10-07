package middleware

import (
	"chi-test/config"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
)

var JwtKey = []byte(config.GetEnv("JWT_SECRET", ""))

func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("JWT authentication middleware called")
		// Get the token from the Authorization header
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			log.Println("Missing token")
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Parse the token and validate it
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil || !token.Valid {
			log.Println("Invalid token:", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		log.Println("Token is valid, proceeding with request")
		// Proceed to the next handler if the token is valid
		next.ServeHTTP(w, r)
	})
}

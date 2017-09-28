package router

import (
	"net/http"
	"strings"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
)

func authMiddleware(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if bearer := r.Header.Get("Authorization"); bearer != "" && len(strings.Split(bearer, " ")) > 1 {
			tokenString := strings.Split(bearer, " ")[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(os.Getenv("authKey")), nil
			})

			if token.Valid {
				fn(w, r)
			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

	}
}

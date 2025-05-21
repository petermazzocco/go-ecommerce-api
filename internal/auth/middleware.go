package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/petermazzocco/go-ecommerce-api/internal/db"
)

func CreateJWT(w http.ResponseWriter, r *http.Request, c db.Cart) (string, error) {
	key := os.Getenv("JWT_KEY")

	claims := jwt.MapClaims{
		"expiresAt": time.Now().Add(24 * time.Hour),
		"cartID":    c.ID, // pass the cart ID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(key))
	if err != nil {
		log.Println("ERROR SIGNING TOKEN: ", err.Error())
		return "", fmt.Errorf("Permission denied")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "dam-nation-shop",
		Value:    ss,
		HttpOnly: true,
		Secure:   false,
		MaxAge:   int(24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
	})

	return ss, nil
}

func ValidateJWT(tokenString string, r *http.Request) error {
	key := os.Getenv("JWT_KEY")

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	
	if err != nil {
		log.Println("ERROR PARSING TOKEN: ", err.Error())
		return fmt.Errorf("Permission denied")

	}

	return nil
}

func CartMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		token, err := r.Cookie("dam-nation-shop")
		if err != nil {
			log.Println("COOKIE ERROR: ", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Permission denied"))
			return
		}

		if err := ValidateJWT(token.Value, r); err != nil {
			log.Println("VALIDATE JWT ERROR: ", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Permission denied"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

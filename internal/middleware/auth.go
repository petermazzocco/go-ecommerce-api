package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/petermazzocco/go-ecommerce-api/internal/cart"
)

func CreateJWT(w http.ResponseWriter, r *http.Request, c *cart.Cart) (string, error) {
	key := os.Getenv("JWT_KEY")

	claims := jwt.MapClaims{
		"expiresAt": time.Now().Add(24 * time.Hour),
		"sessionId": c.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("An unknown error occurred")
	}

    http.SetCookie(w, &http.Cookie{
        Name: "dam-nation-shop",
        Value: ss,
        HttpOnly: true,
		Secure: false,
		MaxAge: int(24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
    })

	return ss, nil
}

func ValidateJWT(tokenString string) error {
	key := os.Getenv("JWT_KEY")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok {
		// Check claims for sessionID of cart
		return nil
	} else {
		return fmt.Errorf("An error occurred")
	}
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		token, err := r.Cookie("dam-nation-shop")
		if err != nil {
			w.Write([]byte("Permission denied"))
			return 
		}
		
		if err := ValidateJWT(token.Value); err != nil {
			w.Write([]byte("Permission denied"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

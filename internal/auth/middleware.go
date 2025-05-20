package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petermazzocco/go-ecommerce-api/internal/db"
	"github.com/petermazzocco/go-ecommerce-api/internal/methods"
)

func CreateJWT(w http.ResponseWriter, r *http.Request, c db.Cart, u db.User) (string, error) {
	key := os.Getenv("JWT_KEY")

	claims := jwt.MapClaims{
		"expiresAt": time.Now().Add(24 * time.Hour),
		"cartID":    c.ID, // pass the cart ID
		"userID":    u.ID, // pass an id if available
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("An unknown error occurred")
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
	cId := r.FormValue("cartID")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		claimsID := claims["cartID"]
		if claimsID == cId {
			return nil
		}

		return fmt.Errorf("An error occurred")
	} else {
		return fmt.Errorf("An error occurred")
	}
}

func ValidateAdminJWT(tokenString string, r *http.Request) error {
	key := os.Getenv("JWT_KEY")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID := claims["userID"]
		if userID != "" {
			user, err := methods.GetUser(int32(1)) 
			if err != nil {
				log.Println(err.Error())
				return err
			}

			notAdmin := user.IsAdmin != pgtype.Bool{Bool: false}
			if notAdmin {
				return fmt.Errorf("An error occurred")
			}
		}
		return fmt.Errorf("An error occurred")
	} else {
		return fmt.Errorf("An error occurred")
	}
}

func CartMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		token, err := r.Cookie("dam-nation-shop")
		if err != nil {
			w.Write([]byte("Permission denied"))
			return
		}

		if err := ValidateJWT(token.Value, r); err != nil {
			w.Write([]byte("Permission denied"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the token from the Authorization header
		token, err := r.Cookie("dam-nation-shop")
		if err != nil {
			w.Write([]byte("Permission denied."))
			return
		}

		if err := ValidateAdminJWT(token.Value, r); err != nil {
			w.Write([]byte("Permission denied."))
			return
		}
		next.ServeHTTP(w, r)
	})
}

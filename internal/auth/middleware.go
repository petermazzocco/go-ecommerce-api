package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petermazzocco/go-ecommerce-api/internal/db"
)

func CreateJWT(w http.ResponseWriter, r *http.Request, c db.Cart) (string, error) {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		log.Println("JWT_KEY environment variable is not set")
		return "", fmt.Errorf("Permission denied")
	}

	cookieName := os.Getenv("COOKIE_NAME")
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
		Name:     cookieName,
		Value:    ss,
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(24 * time.Hour),
		MaxAge:   int(24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
	})

	return ss, nil
}

func CreateAdminJWT(w http.ResponseWriter, r *http.Request, u db.User) (string, error) {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		log.Println("JWT_KEY environment variable is not set")
		return "", fmt.Errorf("Permission denied")
	}
	cookieName := os.Getenv("ADMIN_COOKIE_NAME")
	claims := jwt.MapClaims{
		"expiresAt": time.Now().Add(24 * time.Hour),
		"userID":    u.ID, // pass the user ID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(key))
	if err != nil {
		log.Println("ERROR SIGNING TOKEN: ", err.Error())
		return "", fmt.Errorf("Permission denied")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    ss,
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(24 * time.Hour),
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

func ValidateAdminJWT(tokenString string, r *http.Request) error {
	key := os.Getenv("JWT_KEY")
	url := os.Getenv("DB_URL")

	conn, err := pgx.Connect(r.Context(), url)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close(r.Context())
	q := db.New(conn)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		log.Println("ERROR PARSING TOKEN: ", err.Error())
		return fmt.Errorf("Permission denied")

	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID := claims["userID"]
		if userID != nil {
			id := userID.(float64)
			user, err := q.GetUser(r.Context(), int32(id))
			if err != nil {
				log.Println("ERROR GETTING USER: ", err.Error())
				return fmt.Errorf("Permission denied")
			}
			if user.ID == 0 {
				return fmt.Errorf("Permission denied")
			}
			isAdmin := user.IsAdmin == pgtype.Bool{Bool: true, Valid: true}
			if !isAdmin {
				return fmt.Errorf("Permission denied")
			}
		}
	}
	return nil
}
func CartMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieName := os.Getenv("COOKIE_NAME")
		// Get the token from the Authorization header
		token, err := r.Cookie(cookieName)
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

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieName := os.Getenv("ADMIN_COOKIE_NAME")
		token, err := r.Cookie(cookieName)
		if err != nil {
			log.Println("ADMIN COOKIE ERROR: ", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Permission denied"))
			return
		}

		if err := ValidateAdminJWT(token.Value, r); err != nil {
			log.Println("ADMIN VALIDATE JWT ERROR: ", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Permission denied"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

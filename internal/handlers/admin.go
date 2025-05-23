package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/petermazzocco/go-ecommerce-api/internal/auth"
	"github.com/petermazzocco/go-ecommerce-api/internal/methods"
)

func LoginHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "application/json")

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := methods.Login(ctx, conn, email, password)
	if err != nil {
		log.Println("ERROR LOGGING IN: ", err.Error())
		http.Error(w, "An error occurred authenticating", http.StatusInternalServerError)
		return
	}

	ok, err := methods.CheckUserAdmin(ctx, conn, user.ID)
	if err != nil {
		log.Println("ERROR ADMIN CHECK: ", err.Error())
		http.Error(w, "An error occurred authenticating", http.StatusInternalServerError)
		return
	}

	if ok {
		_, err := auth.CreateAdminJWT(w, r, user)
		if err != nil {
			log.Println("ERROR CREATING JWT: ", err.Error())
			http.Error(w, "An error occurred authenticating", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Logged in!"))
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:    "dam-nation-shop-admin",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out."))
}

func RegisterAdminUserHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "application/json")

	email := r.FormValue("email")
	password := r.FormValue("password")
 
	user, err := methods.CreateUser(ctx, conn, email, password)
	if err != nil {
		log.Println("CREATE USER ERROR: ", err.Error())
		http.Error(w, "An error occurred creating user", http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(user)
	if err != nil {
		log.Println("MARSHAL USER ERROR: ", err.Error())
		http.Error(w, "An error occurred creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

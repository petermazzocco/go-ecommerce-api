package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/petermazzocco/go-ecommerce-api/internal/auth"
	"github.com/petermazzocco/go-ecommerce-api/internal/db"
	"github.com/petermazzocco/go-ecommerce-api/internal/methods"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "application/json")
	
	email := r.FormValue("email")
	passwordHash := r.FormValue("password")
	user, err := methods.CreateUser(ctx, conn, email, passwordHash)
	if err != nil {
		log.Println("CREATE USER ERROR: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var u db.User
	u.ID = user.ID
	u.Email = user.Email
	u.PasswordHash = user.PasswordHash
	_, err = auth.CreateAdminJWT(w, r, u)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(user)
	if err != nil {
		log.Println("MARSHAL USER ERROR: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	strId, _ := strconv.Atoi(id)

	user, err := methods.GetUser(ctx, conn, int32(strId))
	if err != nil {
		log.Println("GET USER ERROR: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(user)
	if err != nil {
		log.Println("MARSHAL USER ERROR: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Response, ctx context.Context, conn *pgx.Conn) {
}

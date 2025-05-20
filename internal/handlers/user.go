package handlers

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Response, ctx context.Context, conn *pgx.Conn) {}

func GetUserHandler(w http.ResponseWriter, r *http.Response, ctx context.Context, conn *pgx.Conn) {}

func DeleteUserHandler(w http.ResponseWriter, r *http.Response, ctx context.Context, conn *pgx.Conn) {}

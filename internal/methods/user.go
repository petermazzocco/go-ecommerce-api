package methods

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petermazzocco/go-ecommerce-api/internal/db"
)

func CreateUser(ctx context.Context, conn *pgx.Conn, email, password string) (db.User, error) {
	q := db.New(conn)

	user, err := q.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: password,
		IsAdmin:      pgtype.Bool{Bool: true, Valid: true},
	})

	if err != nil {
		return db.User{}, err
	}

	return user, nil
}

func CheckUserAdmin(ctx context.Context, id int32) (bool, error) {
	url := os.Getenv("DB_URL")

	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	q := db.New(conn)

	user, err := q.GetUser(ctx, id)
	if err != nil {
		return false, err
	}
	if user.ID == 0 {
		return false, nil
	}
	isAdmin := user.IsAdmin == pgtype.Bool{Bool: true, Valid: true}
	return isAdmin, nil
}

func GetUser(ctx context.Context, conn *pgx.Conn, id int32) (db.User, error) {
	q := db.New(conn)
	user, err := q.GetUser(ctx, id)
	if err != nil {
		return db.User{}, err
	}
	if user.ID == 0 {
		return db.User{}, nil
	}

	return user, nil
}

func DeleteUser(ctx context.Context, conn *pgx.Conn) error {
	return nil
}

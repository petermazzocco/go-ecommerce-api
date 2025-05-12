package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func RunDB(ctx context.Context) (*pgx.Conn, error) {
	url := os.Getenv("DB_URL")
	
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	log.Println("DB running")
	return conn, nil
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)


func main() {
	// Load ENV
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load local env")
	}

	// Start db
	ctx := context.Background()
	url := os.Getenv("DB_URL")

	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB Running")
	defer conn.Close(ctx)
}

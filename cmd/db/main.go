package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	"github.com/petermazzocco/go-ecommerce-api/internal/db"
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

	// create a product in the db so we can test the get product endpoint
	q := db.New(conn)
	// Create a product
	product, err := q.CreateProduct(ctx, db.CreateProductParams{
		Name:        "Test Product",
		Description: pgtype.Text{String: "Test Description", Valid: true},
		Price:       pgtype.Numeric{Int: big.NewInt(19), Valid: true},
	})

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created product: %+v\n", product)

	defer conn.Close(ctx)
}

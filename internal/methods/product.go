package methods

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petermazzocco/go-ecommerce-api/internal/db"
)

type Product struct {
	ID          int      `json:"id"`
	Price       float64  `json:"price"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	Sizes       []Size   `json:"sizes"`
	FitGuide    FitGuide `json:"fitGuide"`
}

type Size struct {
	Size  string `json:"size"`
	Stock int    `json:"stock"`
}

type FitGuide struct {
	BodyLength    float64 `json:"bodyLength"`
	SleeveLength  float64 `json:"sleeveLength"`
	ChestWidth    float64 `json:"chestWidth"`
	ShoulderWidth float64 `json:"shoulderWidth"`
	ArmHole       float64 `json:"armHole"`
	FrontRise     float64 `json:"frontRise"`
	Inseam        float64 `json:"inseam"`
	Hem           float64 `json:"hem"`
	BackRise      float64 `json:"backRise"`
	Waist         float64 `json:"waist"`
	Thigh         float64 `json:"thigh"`
	Knee          float64 `json:"knee"`
}

func CreateProduct(ctx context.Context, conn *pgx.Conn, p Product) (db.Product, error) {
	q := db.New(conn)

	product, err := q.CreateProduct(ctx, db.CreateProductParams{
		Name: p.Name,
		Description: pgtype.Text{String: p.Description},
	})
	
	if err != nil {
		return db.Product{}, fmt.Errorf("Error occurred creating product")
	}
	return product, nil
}

func RemoveProduct(id int) error {
	return nil
}

func AddProductToCollection(pID, cID int) error {
	return nil
}

func RemoveProductFromCollection(pID, cID int) error {
	return nil
}

func UpdateProduct(id int) error {
	return nil
}

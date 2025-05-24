package methods

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petermazzocco/go-ecommerce-api/internal/db"
)

type Product struct {
	ID          int                     `json:"id"`
	Price       float64                 `json:"price"`
	Name        string                  `json:"name"`
	PriceID   string                  `json:"productID"`
	Description string                  `json:"description"`
	Images      []string                `json:"images"`
	Sizes       []db.GetProductSizesRow `json:"sizes"`
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

func GetProducts(ctx context.Context, conn *pgx.Conn) ([]Product, error) {
	var p []Product
	q := db.New(conn)

	products, err := q.ListProducts(ctx)
	for i := range products {
		images, err := q.GetProductImages(ctx, products[i].ID)
		if err != nil {
			log.Println(err.Error())
			return []Product{}, fmt.Errorf("Error occurred fetching product")
		}

		sizes, err := q.GetProductSizes(ctx, products[i].ID)
		if err != nil {
			log.Println(err.Error())
			return []Product{}, fmt.Errorf("Error occurred fetching product")
		}
		floatP, _ := products[i].Price.Float64Value()
		p = append(p, Product{
			ID:          int(products[i].ID),
			Name:        products[i].Name,
			Description: products[i].Description.String,
			Price:       floatP.Float64,
			Images:      images,
			Sizes:       sizes,
		})
	}
	if err != nil {
		log.Println(err.Error())
		return []Product{}, fmt.Errorf("Error occurred fetching product")
	}

	return p, nil
}

func GetProductByID(ctx context.Context, conn *pgx.Conn, id int32) (Product, error) {
	var p Product
	q := db.New(conn)

	product, err := q.GetProduct(ctx, id)
	if err != nil {
		log.Println(err.Error())
		return Product{}, fmt.Errorf("Error occurred fetching product")
	}

	images, err := q.GetProductImages(ctx, id)
	if err != nil {
		log.Println(err.Error())
		return Product{}, fmt.Errorf("Error occurred fetching product")
	}

	sizes, err := q.GetProductSizes(ctx, id)
	if err != nil {
		log.Println(err.Error())
		return Product{}, fmt.Errorf("Error occurred fetching product")
	}
	floatP, _ := product.Price.Float64Value()
	p.ID = int(product.ID)
	p.Name = product.Name
	p.Price = floatP.Float64
	p.Description = product.Description.String
	p.Images = images
	p.Sizes = sizes

	return p, nil
}

func GetProductSizesByID(ctx context.Context, conn *pgx.Conn, id int32) ([]db.GetProductSizesRow, error) {
	q := db.New(conn)

	sizes, err := q.GetProductSizes(ctx, id)

	if err != nil {
		log.Println(err.Error())
		return []db.GetProductSizesRow{}, fmt.Errorf("Error occurred fetching product")

	}

	return sizes, nil
}

func GetProductImagesByID(ctx context.Context, conn *pgx.Conn, id int32) ([]string, error) {
	q := db.New(conn)

	images, err := q.GetProductImages(ctx, id)

	if err != nil {
		log.Println(err.Error())
		return []string{}, fmt.Errorf("Error occurred fetching product")
	}

	return images, nil
}

func CreateProduct(ctx context.Context, conn *pgx.Conn, p Product) (db.Product, error) {
	q := db.New(conn)

	var price pgtype.Numeric
	strPrice := strconv.FormatFloat(p.Price, 'f', -1, 64)
	err := price.Scan(strPrice)
	if err != nil {
		log.Println(err.Error())
		return db.Product{}, fmt.Errorf("Error occurred creating product")
	}
	product, err := q.CreateProduct(ctx, db.CreateProductParams{
		Name:        p.Name,
		Description: pgtype.Text{String: p.Description},
		Price:       price,
		PriceID:   p.PriceID,
	})

	if err != nil {
		log.Println(err.Error())
		return db.Product{}, fmt.Errorf("Error occurred creating product")
	}
	return product, nil
}

func AddProductSizes(ctx context.Context, conn *pgx.Conn, pID int32, sizes []Size) error {
	q := db.New(conn)

	for _, size := range sizes {
		if err := q.AddProductSize(ctx, db.AddProductSizeParams{
			ProductID: pID,
			SizeName:  size.Size,
			Stock:     int32(size.Stock),
		}); err != nil {
			log.Println(err.Error())
			return fmt.Errorf("Error occurred creating product")
		}
	}

	return nil
}

func AddProductImages(ctx context.Context, conn *pgx.Conn, pID int32, images []string) error {
	q := db.New(conn)

	for _, image := range images {
		if err := q.AddProductImage(ctx, db.AddProductImageParams{
			ProductID: pID,
			ImageUrl:  image,
		}); err != nil {
			log.Println(err.Error())
			return fmt.Errorf("Error occurred creating product")
		}
	}

	return nil
}

func RemoveProduct(ctx context.Context, conn *pgx.Conn, id int) error {
	q := db.New(conn)

	_, err := GetProductByID(ctx, conn, int32(id))
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("Product ID is not valid")
	}

	if err := q.DeleteProduct(ctx, int32(id)); err != nil {
		return fmt.Errorf("Error occurred creating product")
	}
	return nil
}

func UpdateProduct(ctx context.Context, conn *pgx.Conn, p Product) error {
	q := db.New(conn)

	_, err := GetProductByID(ctx, conn, int32(p.ID))
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("Product ID is not valid")
	}

	var price pgtype.Numeric

	strPrice := strconv.FormatFloat(p.Price, 'f', -1, 64)
	if err := price.Scan(strPrice); err != nil {
		log.Println(err.Error())
		return fmt.Errorf("Error occurred updating product")
	}

	if err := q.UpdateProduct(ctx, db.UpdateProductParams{
		ID:          int32(p.ID),
		Name:        p.Name,
		Price:       price,
		Description: pgtype.Text{String: p.Description},
	}); err != nil {
		return fmt.Errorf("Error occurred updating product")
	}

	return nil
}



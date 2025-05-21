package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/petermazzocco/go-ecommerce-api/internal/db"
	"github.com/petermazzocco/go-ecommerce-api/internal/methods"
)

type NewProductHandler struct {
	Product db.Product     `json:"product"`
	Images  []string       `json:"images"`
	Sizes   []methods.Size `json:"sizes"`
}

func CreateProductHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "application/json")
	
	var p methods.Product
	name := r.PostFormValue("productName")
	description := r.PostFormValue("productDescription")
	price := r.PostFormValue("productPrice")

	p.Name = name
	p.Description = description
	floatPrice, _ := strconv.ParseFloat(price, 64)
	p.Price = floatPrice

	product, err := methods.CreateProduct(ctx, conn, p)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := methods.AddProductSizes(ctx, conn, product.ID, []methods.Size{
		{Size: "S", Stock: 10},
		{Size: "M", Stock: 20},
		{Size: "L", Stock: 15},
		{Size: "XL", Stock: 5},
	}); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := methods.AddProductImages(ctx, conn, product.ID, []string{
		"https://example.com/image1.jpg",
		"https://example.com/image2.jpg",
	}); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(product)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func ListProductsHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "application/json")

	products, err := methods.GetProducts(ctx, conn)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal the products into JSON
	j, err := json.Marshal(products)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func GetProductHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	strId, _ := strconv.Atoi(id)
	product, err := methods.GetProductByID(ctx, conn, int32(strId))

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(product)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "text/plain")

	id := chi.URLParam(r, "id")
	idInt, _ := strconv.Atoi(id)
	if err := methods.RemoveProduct(ctx, conn, idInt); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product deleted"))
}

func UpdateProductHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "text/plain")

	id := chi.URLParam(r, "id")
	name := r.PostFormValue("productName")
	description := r.PostFormValue("productDescription")
	price := r.PostFormValue("productPrice")

	var p methods.Product
	p.Name = name
	p.Description = description
	floatPrice, _ := strconv.ParseFloat(price, 64)
	p.Price = floatPrice
	p.ID, _ = strconv.Atoi(id)

	if err := methods.UpdateProduct(ctx, conn, p); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product updated"))
}

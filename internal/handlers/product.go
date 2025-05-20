package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/petermazzocco/go-ecommerce-api/internal/methods"
)

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

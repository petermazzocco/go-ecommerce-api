package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petermazzocco/go-ecommerce-api/internal/db"
	"github.com/petermazzocco/go-ecommerce-api/internal/methods"
)

func GetCollectionsHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "application/json")

	collections, err := methods.GetCollections(ctx, conn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(collections) == 0 {
		empty := make([]db.Collection, 0)
		json, err := json.Marshal(empty)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(json)
		return
	}

	json, err := json.Marshal(collections)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func CreateCollectionHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "application/json")

	name := r.FormValue("name")
	description := r.FormValue("description")
	var c db.Collection

	c.Name = name
	c.Description = pgtype.Text{String: description, Valid: true}
	collection, err := methods.CreateCollection(ctx, conn, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(collection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func GetCollectionByIDHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection, err := methods.GetCollection(ctx, conn, intId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(collection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func DeleteCollectionByIDHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "text/plain")

	id := chi.URLParam(r, "id")
	idInt, _ := strconv.Atoi(id)
	if err := methods.DeleteCollection(ctx, conn, idInt); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Collection deleted"))
}

func UpdateCollectionByIDHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "text/plain")

	// id := chi.URLParam(r, "id")
	var c db.Collection
	// idInt,V _ := strconv.Atoi(id)
	if err := methods.UpdateCollection(ctx, conn, c); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Collection updated"))
}

func AddProductToCollectionHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "text/plain")

	collectionID := chi.URLParam(r, "collectionID")
	productID := chi.URLParam(r, "productID")
	collectionIDInt, _ := strconv.Atoi(collectionID)
	productIDInt, _ := strconv.Atoi(productID)

	if err := methods.AddProductToCollection(ctx, conn, collectionIDInt, productIDInt); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product added to collection"))
}

func RemoveProductFromCollectionHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "text/plain")

	collectionID := chi.URLParam(r, "collectionID")
	productID := chi.URLParam(r, "productID")
	collectionIDInt, _ := strconv.Atoi(collectionID)
	productIDInt, _ := strconv.Atoi(productID)

	if err := methods.RemoveProductFromCollection(ctx, conn, collectionIDInt, productIDInt); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product removed from collection"))
}

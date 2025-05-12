package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/petermazzocco/go-ecommerce-api/internal/methods"
	"github.com/petermazzocco/go-ecommerce-api/internal/auth"
)

func GetCookie(r *http.Request) string {
	cookie, err := r.Cookie("dam_nation_shop")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func NewCartHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	// Create a new cart pointer
	cart := methods.NewCart(ctx, conn)

	// Create the JWT when we create a new NewCartHandler
	_, err := auth.CreateJWT(w, r, cart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("New cart created: " + cart.ID.String()))
}

func GetCartProductsHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "application/json")
	id := r.FormValue("cartID")
	p, _ := uuid.Parse(id)
	items, err := methods.GetItems(ctx, conn, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

func ClearCartHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "application/json")
	id := r.FormValue("cartID")
	p, _ := uuid.Parse(id)
	if err := methods.ClearAll(ctx, conn, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:    "dam_nation_shop",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1, // or 0
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cart has been cleared"))
}

func AddItemHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	c := r.PostFormValue("cartID")
	p, _ := uuid.Parse(c)

	prod := r.PostFormValue("productID")
	prodID, err := strconv.Atoi(prod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := methods.AddItem(ctx, conn, p, prodID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item has been added to cart"))
}

func RemoveItemHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	c := r.PostFormValue("cartID")
	p, _ := uuid.Parse(c)

	prod := r.PostFormValue("productID")
	prodID, err := strconv.Atoi(prod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := methods.RemoveItem(ctx, conn, p, prodID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item has been removed from cart"))
}

func UpdateItemQuantityHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {

	c := r.PostFormValue("cartID")
	q := r.PostFormValue("quantity")
	p, _ := uuid.Parse(c)

	quan, err := strconv.Atoi(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prod := r.PostFormValue("productID")
	prodID, err := strconv.Atoi(prod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := methods.UpdateItemQuantity(ctx, conn, p, prodID, quan); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item has been removed from cart"))
}

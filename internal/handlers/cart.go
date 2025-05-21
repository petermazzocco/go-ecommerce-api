package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/petermazzocco/go-ecommerce-api/internal/auth"
	"github.com/petermazzocco/go-ecommerce-api/internal/methods"
)

func GetCartIDFromCookie(r *http.Request) (string, error) {
	key := os.Getenv("JWT_KEY")

	cookie, err := r.Cookie("dam-nation-shop")
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		log.Println(err.Error())
		return "", fmt.Errorf("An unknown error occurred")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		cartID := claims["cartID"]
		if cartID != nil {
			return cartID.(string), nil
		}
	}

	return "", fmt.Errorf("An unknown error occurred")
}

func NewCartHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "text/plain")

	cart, err := methods.NewCart(ctx, conn)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	_, err = auth.CreateJWT(w, r, cart)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("New cart created: " + cart.ID.String()))
}

func GetCartProductsHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "application/json")

	id, err := GetCartIDFromCookie(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	p, err := uuid.Parse(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	items, err := methods.GetItems(ctx, conn, p)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(items)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

func ClearCartHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "text/plain")

	id, err := GetCartIDFromCookie(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	p, err := uuid.Parse(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	if err := methods.ClearAll(ctx, conn, p); err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
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
	w.Header().Set("Content-Type", "text/plain")

	id, err := GetCartIDFromCookie(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	p, err := uuid.Parse(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	prod := r.PostFormValue("productID")
	prodID, err := strconv.Atoi(prod)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	quan := r.PostFormValue("quantity")
	q, err := strconv.Atoi(quan)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}
	_, err = methods.GetProductByID(ctx, conn, int32(prodID))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	if err := methods.AddItem(ctx, conn, p, prodID, q); err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item has been added to cart"))
}

func RemoveItemHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "text/plain")

	id, err := GetCartIDFromCookie(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	p, err := uuid.Parse(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	prod := r.PostFormValue("productID")
	prodID, err := strconv.Atoi(prod)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	if err := methods.RemoveItem(ctx, conn, p, prodID); err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item has been removed from cart"))
}

func UpdateItemQuantityHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "text/plain")

	id, err := GetCartIDFromCookie(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	q := r.PostFormValue("quantity")
	quan, err := strconv.Atoi(q)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	prod := r.PostFormValue("productID")
	prodID, err := strconv.Atoi(prod)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	p, err := uuid.Parse(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	if err := methods.UpdateItemQuantity(ctx, conn, p, prodID, quan); err != nil {
		log.Println(err.Error())
		http.Error(w, "An unknown error occurred", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item has been updated in the cart"))
}

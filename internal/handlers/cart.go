package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/petermazzocco/go-ecommerce-api/internal/cart"
)

func GetCookie(r *http.Request) string {
	cookie, err := r.Cookie("dam_nation_shop")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func GetCartProductsHandler(w http.ResponseWriter, r *http.Request, c *cart.Cart) {
	w.Header().Set("Content-Type", "application/json")
	cookie := GetCookie(r)
	items, err := c.GetItems(cookie)
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

func ClearCartHandler(w http.ResponseWriter, c *cart.Cart) {
	if err := c.ClearAll(); err != nil {
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

func NewCartHandler(w http.ResponseWriter, r *http.Request) {
	cart := cart.NewCart()

	// Create a new cookie
	cookie := &http.Cookie{
		Name:     "dam_nation_shop",
		Value:    "12345",
		Path:     "/",                            // Adjust path as needed
		Expires:  time.Now().Add(24 * time.Hour), // Set expiration time
		HttpOnly: true,                           // Optional: prevents client-side JavaScript access
		Secure:   true,                           // Optional: requires HTTPS
		SameSite: http.SameSiteStrictMode,        // Optional: helps prevent CSRF attacks
	}

	http.SetCookie(w, cookie)
	fmt.Println("Cart", cart)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("New cart has been creatd with the cookie"))

}

func AddItemHandler(w http.ResponseWriter, r *http.Request, c *cart.Cart) {
	item := r.PostFormValue("itemID")
	id, err := strconv.Atoi(item)
	if err != nil {
		w.Write([]byte("An error occurred formatting ID to int"))
		return
	}

	if err := c.AddItem(id); err != nil {
		formatted := fmt.Sprintf("An error occurred: %s", err.Error())
		w.Write([]byte(formatted))
		return
	}

	w.Write([]byte("Item has been added to the cart"))
}

func RemoveItemHandler(w http.ResponseWriter, r *http.Request, c *cart.Cart) {
	item := r.PostFormValue("itemID")
	id, err := strconv.Atoi(item)
	if err != nil {
		w.Write([]byte("An error occurred formatting ID to int"))
		return
	}


	cookie := GetCookie(r)
	if err := c.RemoveItem(id, cookie); err != nil {
		formatted := fmt.Sprintf("An error occurred: %s", err.Error())
		w.Write([]byte(formatted))
		return
	}
}

func RemoveItemQuantityHandler(w http.ResponseWriter, r *http.Request) {

}

//
// func AddItemQuantityHandler(w http.ResponseWriter, r *http.Request) {
// 	cart := &cart.Cart{}
//
// 	item := r.PostFormValue("itemID")
// 	id, err := strconv.Atoi(item)
// 	if err != nil {
// 		w.Write([]byte("An error occurred formatting ID to int"))
// 		return
// 	}
//
// 	c, err := cart.AddItemQuantity(id)
// 	if err != nil {
// 		formatted := fmt.Sprintf("An error occurred: %s", err.Error())
// 		w.Write([]byte(formatted))
// 		return
// 	}
//
// 	formatted := fmt.Sprintf("Quantity for item %v has been added: %v", id, c)
// 	w.Write([]byte(formatted))
// }

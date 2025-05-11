package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/petermazzocco/go-ecommerce-api/internal/cart"
	"github.com/petermazzocco/go-ecommerce-api/internal/handlers"
	auth "github.com/petermazzocco/go-ecommerce-api/internal/middleware"
)

func main() {
	cart := cart.NewCart()
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load local env")
	}
	// available api routes for our backend
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Api route
	r.Route("/api", func(r chi.Router) {

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("API"))
		})

		// products route group
		r.Route("/products", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("All products here. Incorporate pagination for all products"))
			})
			r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
				id := chi.URLParam(r, "id")
				w.Write([]byte(id))
			})
		})

		r.Route("/collections", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("All collections here. Incorporate pagination for all collections"))
			})
			r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
				id := chi.URLParam(r, "id")
				w.Write([]byte(id))
			})
		})

		// Any time we use the cart, we need the JWT and session ID of the cart to track
		r.Post("/new-cart", func(w http.ResponseWriter, r *http.Request) {
			// Create the JWT when we create a new cart
			_, err := auth.CreateJWT(w, r, cart)
			// Set token as a cookie or in the response body
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error"))
			}


			handlers.NewCartHandler(w, r)
		})
		// Cart logic ( must incorporate jwt to persist cart items over time )
		r.Route("/cart", func(r chi.Router) {
			r.Use(auth.Middleware)
			r.Get("/products", func(w http.ResponseWriter, r *http.Request) {
				handlers.GetCartProductsHandler(w, r, cart)
			})
			r.Post("/clear", func(w http.ResponseWriter, r *http.Request) {
				handlers.ClearCartHandler(w, cart)
			})
			r.Post("/add-item", func(w http.ResponseWriter, r *http.Request) {
				handlers.AddItemHandler(w, r, cart)
			})
			r.Post("/remove-item", func(w http.ResponseWriter, r *http.Request) {
				handlers.RemoveItemHandler(w, r, cart)
			})
			// r.Post("/add-quantity-item", handlers.AddItemQuantity)
			// r.Post("/remove-quantity-item", handlers.RemoveItemQuantity)
		})

	})

	http.ListenAndServe(":8080", r)
}

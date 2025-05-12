package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/petermazzocco/go-ecommerce-api/internal/auth"
	"github.com/petermazzocco/go-ecommerce-api/internal/db"
	"github.com/petermazzocco/go-ecommerce-api/internal/handlers"
)

func main() {
	// Load ENV
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load local env")
	}

	// Start db 	
	ctx := context.Background()
	conn, err := db.RunDB(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Chi routers
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Api route
	r.Route("/api", func(r chi.Router) {
		// Status check
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Ecommerce API"))
		})

		// Public facing products route group to return product information
		r.Route("/products", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("All products here. Incorporate pagination for all products"))
			})
			r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
				id := chi.URLParam(r, "id")
				w.Write([]byte(id))
			})
		})

		// Public facing collections route group to return collection information
		r.Route("/collections", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("All collections here. Incorporate pagination for all collections"))
			})
			r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
				id := chi.URLParam(r, "id")
				w.Write([]byte(id))
			})
		})

		// `new-cart` will create a new cart, JWT and session storage
		r.Post("/new-cart", func(w http.ResponseWriter, r *http.Request) {
			handlers.NewCartHandler(w, r, ctx, conn)
		})

		// Cart logic ( must incorporate jwt to persist cart items over time )
		r.Route("/cart", func(r chi.Router) {
			r.Use(auth.Middleware) // Require each route has a valid JWT and cart session ID
			r.Get("/products", func(w http.ResponseWriter, r *http.Request) {
				handlers.GetCartProductsHandler(w, r, ctx, conn)
			})
			r.Post("/clear", func(w http.ResponseWriter, r *http.Request) {
				handlers.ClearCartHandler(w, r, ctx, conn)
			})
			r.Post("/add-item", func(w http.ResponseWriter, r *http.Request) {
				handlers.AddItemHandler(w, r, ctx, conn)
			})
			r.Post("/remove-item", func(w http.ResponseWriter, r *http.Request) {
				handlers.RemoveItemHandler(w, r, ctx, conn)
			})
			r.Post("/update-quantity-item", func(w http.ResponseWriter, r *http.Request) {
				handlers.UpdateItemQuantityHandler(w, r, ctx, conn)
			})	
		})

	})

	http.ListenAndServe(":8080", r)
}

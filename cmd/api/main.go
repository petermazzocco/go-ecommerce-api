package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/petermazzocco/go-ecommerce-api/internal/auth"
	"github.com/petermazzocco/go-ecommerce-api/internal/handlers"
)

func main() {
	// Load ENV
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load local env")
	}

	// Start db
	ctx := context.Background()
	url := os.Getenv("DB_URL")

	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	// Chi routers
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		// Status check
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Ecommerce API"))
		})

		r.Route("/admin", func(r chi.Router) {
			r.Route("/products", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					handlers.ListProductsHandler(w, r, ctx, conn)
				})
				// Get specific product
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", func(w http.ResponseWriter, r *http.Request) {
						handlers.GetProductHandler(w, r, ctx, conn)
					})
					r.Put("/", func(w http.ResponseWriter, r *http.Request) {
						handlers.UpdateProductHandler(w, r, ctx, conn)
					})
					r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
						handlers.DeleteProductHandler(w, r, ctx, conn)
					})
				})
				// Create new product
				r.Post("/new", func(w http.ResponseWriter, r *http.Request) {
					handlers.CreateProductHandler(w, r, ctx, conn)
				})
			})

			// Public facing collections route group to return collection information
			r.Route("/collections", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					handlers.GetCollectionsHandler(w, r, ctx, conn)
				})
				// Get specific collection
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", func(w http.ResponseWriter, r *http.Request) {
						handlers.GetCollectionByIDHandler(w, r, ctx, conn)
					})
					r.Put("/", func(w http.ResponseWriter, r *http.Request) {
						handlers.UpdateCollectionByIDHandler(w, r, ctx, conn)
					})
					r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
						handlers.DeleteCollectionByIDHandler(w, r, ctx, conn)
					})
				})
				// Create new collection
				r.Post("/new", func(w http.ResponseWriter, r *http.Request) {
					handlers.CreateCollectionHandler(w, r, ctx, conn)
				})
			})
		})

		// Public facing products route group to return product information
		r.Route("/products", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				handlers.ListProductsHandler(w, r, ctx, conn)
			})
			r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
				handlers.GetProductHandler(w, r, ctx, conn)
			})
		})

		// Public facing collections route group to return collection information
		r.Route("/collections", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				handlers.GetCollectionsHandler(w, r, ctx, conn)
			})
			r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
				handlers.GetCollectionByIDHandler(w, r, ctx, conn)
			})
		})

		// `new-cart` will create a new cart, JWT 
		r.Post("/new-cart", func(w http.ResponseWriter, r *http.Request) {
			handlers.NewCartHandler(w, r, ctx, conn)
		})

		// Cart logic ( must incorporate jwt to persist cart items over time )
		r.Route("/cart", func(r chi.Router) {
			r.Use(auth.CartMiddleware) // Require each route has a valid JWT and cart session ID
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

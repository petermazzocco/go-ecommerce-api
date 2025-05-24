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
		// Health check
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Ecommerce API"))
		})

		// Portal login for admin users
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
				handlers.LoginHandler(w, r, ctx, conn)
			})
			r.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
				handlers.LogoutHandler(w, r)
			})
			r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
				handlers.RegisterAdminUserHandler(w, r, ctx, conn)
			})
		})

		// Admin route group to require admin role
		r.Route("/admin", func(r chi.Router) {
			r.Use(auth.AdminMiddleware) // Require each route has a valid JWT with a maps claim to an admin role
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Admin Portal"))
			})
			r.Route("/users", func(r chi.Router) {
				r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
					handlers.RegisterAdminUserHandler(w, r, ctx, conn)
				})
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", func(w http.ResponseWriter, r *http.Request) {
						handlers.GetUserHandler(w, r, ctx, conn)
					})
					r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
						handlers.DeleteUserHandler(w, r, ctx, conn)
					})
				})
			})

			// Products route group to manage products
			r.Route("/products", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					handlers.ListProductsHandler(w, r, ctx, conn)
				})
				r.Post("/", func(w http.ResponseWriter, r *http.Request) {
					handlers.CreateProductHandler(w, r, ctx, conn)
				})
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
			})

			// Collections route group to manage collections
			r.Route("/collections", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					handlers.GetCollectionsHandler(w, r, ctx, conn)
				})
				r.Post("/", func(w http.ResponseWriter, r *http.Request) {
					handlers.CreateCollectionHandler(w, r, ctx, conn)
				})
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
					r.Route("/product", func(r chi.Router) {
						// Add or remove products from a collection
						r.Route("/{id}", func(r chi.Router) {
							r.Post("/", func(w http.ResponseWriter, r *http.Request) {
								handlers.AddProductToCollectionHandler(w, r, ctx, conn)
							})
							r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
								handlers.RemoveProductFromCollectionHandler(w, r, ctx, conn)
							})
						})
					})
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

		// Creates a new cart with a unique ID that is stored in a cookie with a JWT for authentication
		r.Post("/new-cart", func(w http.ResponseWriter, r *http.Request) {
			handlers.NewCartHandler(w, r, ctx, conn)
		})

		// Cart route group requires a valid JWT and cart session ID
		r.Route("/cart", func(r chi.Router) {
			r.Use(auth.CartMiddleware) // Require each route has a valid JWT and cart session ID
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				handlers.GetCartProductsHandler(w, r, ctx, conn)
			})
			r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				handlers.ClearCartHandler(w, r, ctx, conn)
			})
			r.Post("/add", func(w http.ResponseWriter, r *http.Request) {
				handlers.AddItemHandler(w, r, ctx, conn)
			})
			// Product ID in the cart to update quan or remove
			r.Route("/{productID}", func(r chi.Router) {
				r.Put("/", func(w http.ResponseWriter, r *http.Request) {
					handlers.UpdateItemQuantityHandler(w, r, ctx, conn)
				})
				r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
					handlers.RemoveItemHandler(w, r, ctx, conn)
				})
			})
			// Create a Stripe check out session
			r.Post("/checkout", func(w http.ResponseWriter, r *http.Request) {
				handlers.CreateCheckoutSession(w, r, ctx, conn)
			})
		})

	})

	http.ListenAndServe(":8080", r)
}

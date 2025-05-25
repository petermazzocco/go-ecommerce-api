package handlers

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/petermazzocco/go-ecommerce-api/internal/methods"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

func CreateCheckoutSession(w http.ResponseWriter, r *http.Request, ctx context.Context, conn *pgx.Conn) {
	w.Header().Set("Content-Type", "text/plain")

	// Load Stripe key from environment variable
	stripe.Key = os.Getenv("STRIPE_KEY")

	// Get cart ID from cookie
	id, err := GetCartIDFromCookie(r)
	if err != nil {
		log.Println("Error getting cart ID from cookie:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the cart ID from the cookie
	strID, err := uuid.Parse(id)
	if err != nil {
		log.Println("Error parsing cart ID:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get cart items from the database
	cart, err := methods.GetCart(ctx, conn, strID)
	if err != nil {
		log.Println("Error getting cart:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get items in the cart
	items, err := methods.GetItems(ctx, conn, strID)
	if err != nil {
		log.Println("Error getting items:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return	
	}

	if len(items) == 0 {
		log.Println("No items in cart")
		http.Error(w, "No items in cart", http.StatusBadRequest)
		return
	}

	// Stripe line items
	var lineItems []*stripe.CheckoutSessionLineItemParams
	for _, item := range items {
		lineItem := &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(string(item.PriceID)), 
			Quantity: stripe.Int64(int64(item.Quantity)),
		}
		lineItems = append(lineItems, lineItem)
	}

	// Create a stripe Customer if needed
	email := r.PostFormValue("email")
	name := r.PostFormValue("name")
	customerParams := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
		// Add more customer details as needed
	}

	customer, err := customer.New(customerParams)

	// Stripe checkout session params
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("https://example.com/success"),
		CancelURL:  stripe.String("https://example.com/cancel"),
		LineItems:  lineItems,
		Metadata: map[string]string{
			"cartID": cart.ID.String(),
			// Add more metadata as needed here
		},
		BillingAddressCollection: stripe.String(stripe.CheckoutSessionBillingAddressCollectionRequired),
		ShippingAddressCollection: &stripe.CheckoutSessionShippingAddressCollectionParams{
			AllowedCountries: stripe.StringSlice([]string{"US", "CA"}), // Example allowed AllowedCountries
		},
		Customer: stripe.String(customer.ID),
		Mode: stripe.String(stripe.CheckoutSessionModePayment),
	}

	// Create the checkout session
	result, err := session.New(params)
	if err != nil {
		log.Println("Error creating checkout session:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result.URL))
	// Redirect to the checkout session url
	http.Redirect(w, r, result.URL, http.StatusSeeOther)
}

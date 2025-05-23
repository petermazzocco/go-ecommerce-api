package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

func CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	stripe.Key = os.Getenv("STRIPE_KEY")
	params := &stripe.CheckoutSessionParams{
	  SuccessURL: stripe.String("https://example.com/success"),
	  LineItems: []*stripe.CheckoutSessionLineItemParams{
	    &stripe.CheckoutSessionLineItemParams{
	      Price: stripe.String("price_1RRwYxA1cCGZeIWJ10NPnT1u"),
	      Quantity: stripe.Int64(2),
	    },
	  },
		Metadata: map[string]string{
			"cart_id": "6969",
		},
	  Mode: stripe.String(stripe.CheckoutSessionModePayment),
	};
	result, err := session.New(params);
	if err != nil {
		log.Println("Error creating checkout session:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Creating checkout session..." + result.URL))
	// Redirect to the checkout session url
	http.Redirect(w, r, result.URL, http.StatusSeeOther)
}

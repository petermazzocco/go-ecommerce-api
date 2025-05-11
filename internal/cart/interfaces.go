package cart

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/petermazzocco/go-ecommerce-api/internal/product"
)

type Cart struct {
	ID uuid.UUID
	Products []CartProduct
}

type CartProduct struct {
	Product  product.Product
	Quantity int
}

type ShoppingCart interface {
	GetCart(string) (*Cart, error)
	GetItems(string) ([]CartProduct, error)
	NewCart() *Cart
	AddItem(int, string) error
	ClearAll() error
	RemoveItem(int, string) error
	RemoveItemQuantity(int) error
	AddItemQuantity(int) error
}

func (c *Cart) GetCart(cookie string) (*Cart, error) {
	if cookie == "" {
		return nil, fmt.Errorf("Cookie invalid")
	}
	return c, nil
}

// Get all current products in a cart
func (c *Cart) GetItems(cookie string) ([]CartProduct, error) {
	if cookie == "" {
		return nil, fmt.Errorf("Cookie invalid")
	}
	return c.Products, nil
}

// Create a new cart
func NewCart() *Cart {
	return &Cart{
		ID: uuid.New(),
		Products: make([]CartProduct, 0),
	}
}

// Clear all cart items
func (c *Cart) ClearAll() error {
	if len(c.Products) == 0 {
		return fmt.Errorf("No items in cart")
	}
	c.Products = make([]CartProduct, 0)
	return nil
}

// Remove an item from cart
func (c *Cart) RemoveItem(id int, cookie string) error {
	// Loop through the current cart products
	for i, p := range c.Products {
		// If the ID's match, delete it
		if p.Product.ID == id {
			// Remove the current product index from the array of cart.Products
			c.Products = slices.Delete(c.Products, i, i+1)
		} else {
			return fmt.Errorf("This item is not in your cart")
		}
	}

	return nil
}

// Add items on to existing cart
func (c *Cart) AddItem(id int) error {
	// Make a cehck if the product ID exists in our db
	if id == 0 {
		return fmt.Errorf("Invalid product ID")
	}

	// Check if the product exists in the cart
	for _, p := range c.Products {
		// If the product already exists in the Cart
		// increment the quantity
		if p.Product.ID == id {
			if err := c.AddItemQuantity(id); err != nil {
				return fmt.Errorf("Error occurred incrementing item")
			}
			return nil
		} else {
			// Add it as a new item in the cart.Products
			c.Products = append(c.Products, p)
		}
	}
	return nil
}

// Remove quanity from item in cart
func (c *Cart) RemoveItemQuantity(id int) error {
	// Make a cehck if the product ID exists in our db
	if id <= 0 {
		return fmt.Errorf("Invalid product ID")
	}

	for _, p := range c.Products {
		if p.Product.ID == id {
			p.Quantity--
		} else {
			return fmt.Errorf("Cannot not decrement an item that is not in cart")
		}
	}
	return nil
}

// Add quantity to item in cart
func (c *Cart) AddItemQuantity(id int) error {
	// Make a cehck if the product ID exists in our db
	if id == 0 {
		return fmt.Errorf("Invalid product ID")
	}

	for _, p := range c.Products {
		if p.Product.ID == id {
			p.Quantity++
		} else {
			return fmt.Errorf("Cannot not increment an item that is not in cart")
		}
	}
	return nil
}

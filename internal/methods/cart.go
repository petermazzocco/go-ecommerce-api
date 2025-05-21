package methods

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petermazzocco/go-ecommerce-api/internal/db"
)

func NewCart(ctx context.Context, conn *pgx.Conn) (db.Cart, error) {
	q := db.New(conn)

	cart, err := q.CreateCart(ctx, pgtype.UUID{Bytes: uuid.New(), Valid: true})

	if err != nil {
		log.Println("NEW CART ERROR: ", err.Error())
		return db.Cart{}, err
	}

	return cart, nil
}

func GetCart(ctx context.Context, conn *pgx.Conn, id uuid.UUID) (db.Cart, error) {
	q := db.New(conn)

	parsedID := pgtype.UUID{Bytes: id, Valid: true}

	cart, err := q.GetCart(ctx, parsedID)
	if err != nil || cart.ID != parsedID {
		log.Println("GET CART ERROR: ", err.Error())
		return db.Cart{}, fmt.Errorf("Error getting cart")
	}
	return cart, nil
}

func GetItems(ctx context.Context, conn *pgx.Conn, id uuid.UUID) ([]db.GetCartItemsRow, error) {
	q := db.New(conn)

	_, err := GetCart(ctx, conn, id)
	if err != nil {
		log.Println("GET CART ERROR: ", err.Error())
		return nil, fmt.Errorf("Error getting cart")
	}

	items, err := q.GetCartItems(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		log.Println("GET CART ITEMS ERROR: ", err.Error())
		return []db.GetCartItemsRow{}, fmt.Errorf("Error fetching items")
	}

	if len(items) == 0 {
		// return an empty array instead of null for better front end handling
		return make([]db.GetCartItemsRow, 0), nil
	}
	return items, nil
}

func ClearAll(ctx context.Context, conn *pgx.Conn, id uuid.UUID) error {
	q := db.New(conn)

	_, err := GetCart(ctx, conn, id)
	if err != nil {
		log.Println("GET CART ERROR: ", err.Error())
		return fmt.Errorf("Error getting cart")
	}

	if err := q.ClearCart(ctx, pgtype.UUID{Bytes: id, Valid: true}); err != nil {
		log.Println("CLEAR CART ERROR: ", err.Error())
		return fmt.Errorf("Error clearing items in cart ")
	}

	return nil
}

func RemoveItem(ctx context.Context, conn *pgx.Conn, id uuid.UUID, prodID int) error {
	q := db.New(conn)

	_, err := GetCart(ctx, conn, id)
	if err != nil {
		log.Println("GET CART ERROR: ", err.Error())
		return fmt.Errorf("Error getting cart")
	}

	if err := q.RemoveCartItem(ctx, db.RemoveCartItemParams{
		CartID:    pgtype.UUID{Bytes: id, Valid: true},
		ProductID: int32(prodID),
	}); err != nil {
		log.Println("REMOVE CART ITEM ERROR: ", err.Error())
		return fmt.Errorf("Error removing the item in cart")
	}

	return nil
}

func AddItem(ctx context.Context, conn *pgx.Conn, id uuid.UUID, prodID int, quan int) error {
	q := db.New(conn)

	_, err := GetCart(ctx, conn, id)
	if err != nil {
		log.Println("GET CART ERROR: ", err.Error())
		return fmt.Errorf("Error getting cart")
	}

	if err := q.AddCartItem(ctx, db.AddCartItemParams{
		CartID:    pgtype.UUID{Bytes: id, Valid: true},
		ProductID: int32(prodID),
		Quantity:  int32(quan),
	}); err != nil {
		log.Println("ADD CART ITEM ERROR: ", err.Error())
		return fmt.Errorf("Error adding the item in cart")
	}

	return nil
}

func UpdateItemQuantity(ctx context.Context, conn *pgx.Conn, id uuid.UUID, prodID int, quan int) error {
	q := db.New(conn)

	_, err := GetCart(ctx, conn, id)
	if err != nil {
		log.Println("GET CART ERROR: ", err.Error())
		return fmt.Errorf("Error getting cart")
	}

	if err := q.UpdateCartItemQuantity(ctx, db.UpdateCartItemQuantityParams{
		CartID:    pgtype.UUID{Bytes: id, Valid: true},
		ProductID: int32(prodID),
		Quantity:  int32(quan),
	}); err != nil {
		log.Println("UPDATE ITEM CART ERROR :", err.Error())
		return fmt.Errorf("Error changing the item quantity")
	}
	return nil
}

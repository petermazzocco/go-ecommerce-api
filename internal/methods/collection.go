package methods

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/petermazzocco/go-ecommerce-api/internal/db"
)


func CreateCollection(ctx context.Context, conn *pgx.Conn, c db.Collection) (db.Collection, error) {
	q := db.New(conn)

	collection, err := q.CreateCollection(ctx, db.CreateCollectionParams{
		Name:        c.Name,
		Description: c.Description,
	})
	if err != nil {
		log.Println(err.Error())
		return db.Collection{}, err
	}

	return collection, nil
}

func AddProductToCollection(ctx context.Context, conn *pgx.Conn, collectionID int, productID int) error {
	q := db.New(conn)

	if err := q.AddProductToCollection(ctx, db.AddProductToCollectionParams{
		CollectionID: int32(collectionID),
		ProductID:   int32(productID),
	}); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func RemoveProductFromCollection(ctx context.Context, conn *pgx.Conn, collectionID int, productID int) error {
	q := db.New(conn)

	if err := q.RemoveProductFromCollection(ctx, db.RemoveProductFromCollectionParams{
		CollectionID: int32(collectionID),
		ProductID:   int32(productID),
	}); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func GetCollections(ctx context.Context, conn *pgx.Conn) ([]db.Collection, error) {
	q := db.New(conn)

	collections, err := q.ListCollections(ctx)
	if err != nil {
		log.Println(err.Error())
		return []db.Collection{}, err
	}

	return collections, nil
}

func GetCollection(ctx context.Context, conn *pgx.Conn, id int) (db.Collection, error) {
	q := db.New(conn)

	collection, err := q.GetCollection(ctx, int32(id))
	if err != nil {
		log.Println(err.Error())
		return db.Collection{}, err
	}

	return collection, nil
}

func UpdateCollection(ctx context.Context, conn *pgx.Conn, c db.Collection) error {
	q := db.New(conn)

	if err := q.UpdateCollection(ctx, db.UpdateCollectionParams{
		ID:          int32(c.ID),
		Name:        c.Name,
		Description: c.Description,
	});err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func DeleteCollection(ctx context.Context, conn *pgx.Conn, id int) error {
	q := db.New(conn)

	if err := q.DeleteCollection(ctx, int32(id)); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

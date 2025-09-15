package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Order struct {
	ID       uuid.UUID `json:"id"`
	Product  string    `json:"product"`
	Quantity int       `json:"quantity"`
}

type OrderCreatedEvent struct {
	OrderID uuid.UUID `json:"order_id"`
	Product string    `json:"product"`
}

// createOrderInTx creates an order and its corresponding outbox event atomically.
func createOrderInTx(ctx context.Context, tx pgx.Tx, order Order) error {
	// 1. Insert the order
	_, err := tx.Exec(ctx, "INSERT INTO orders (id, product, quantity) VALUES ($1, $2, $3)",
		order.ID, order.Product, order.Quantity)
	if err != nil {
		return err
	}
	log.Printf("Inserted order %s into database", order.ID)

	// 2. Prepare the event message for the outbox
	event := OrderCreatedEvent{
		OrderID: order.ID,
		Product: order.Product,
	}
	msg, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// 3. Insert the event into the outbox table
	_, err = tx.Exec(ctx, "INSERT INTO outbox (topic, message) VALUES ($1, $2, $3)",
		"orders.created", msg)
	if err != nil {
		return err
	}
	log.Printf("Inserted outbox event for order %s", order.ID)

	return nil
}

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Fatalf("Unable to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	newOrder := Order{
		ID:       uuid.New(),
		Product:  "Super Widget",
		Quantity: 10,
	}

	if err := createOrderInTx(ctx, tx, newOrder); err != nil {
		log.Fatalf("Failed to create order: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	log.Println("Successfully created order and outbox event.")
}

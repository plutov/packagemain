package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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

func main() {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	order := Order{
		ID:       uuid.New(),
		Product:  "Super Widget",
		Quantity: 10,
	}

	// Insert order
	_, err = pool.Exec(ctx, "INSERT INTO orders (id, product, quantity) VALUES ($1, $2, $3)",
		order.ID, order.Product, order.Quantity)
	if err != nil {
		log.Fatalf("Failed to insert order: %v", err)
	}
	log.Printf("Inserted order %s", order.ID)

	event := OrderCreatedEvent{OrderID: order.ID, Product: order.Product}
	msg, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Failed to marshal event: %v", err)
	}

	// Publish event
	if err := rdb.Publish(ctx, "orders.created", msg).Err(); err != nil {
		log.Fatalf("Failed to publish event: %v", err)
	}
	log.Printf("Published event for order %s", order.ID)
}

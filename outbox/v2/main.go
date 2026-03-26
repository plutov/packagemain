package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

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

type OutboxMessage struct {
	ID      uuid.UUID
	Topic   string
	Message []byte
}

func createOrder(ctx context.Context, pool *pgxpool.Pool, order Order) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "INSERT INTO orders (id, product, quantity) VALUES ($1, $2, $3)",
		order.ID, order.Product, order.Quantity)
	if err != nil {
		return err
	}

	event := OrderCreatedEvent{OrderID: order.ID, Product: order.Product}
	msg, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "INSERT INTO outbox (topic, message) VALUES ($1, $2)",
		"orders.created", msg)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func relay(ctx context.Context, pool *pgxpool.Pool, rdb *redis.Client) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
		SELECT id, topic, message
		FROM outbox
		WHERE state = 'pending'
		ORDER BY created_at
		LIMIT 1
		FOR UPDATE SKIP LOCKED
	`)
	if err != nil {
		return err
	}

	var msg OutboxMessage
	if !rows.Next() {
		rows.Close()
		return nil
	}
	if err := rows.Scan(&msg.ID, &msg.Topic, &msg.Message); err != nil {
		return err
	}
	rows.Close()

	if err := rdb.Publish(ctx, msg.Topic, msg.Message).Err(); err != nil {
		return err
	}
	log.Printf("Published message %s to topic %s", msg.ID, msg.Topic)

	_, err = tx.Exec(ctx, "UPDATE outbox SET state = 'processed', processed_at = now() WHERE id = $1", msg.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
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

	if err := createOrder(ctx, pool, order); err != nil {
		log.Fatalf("Failed to create order: %v", err)
	}
	log.Printf("Created order %s", order.ID)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if err := relay(context.Background(), pool, rdb); err != nil {
			log.Printf("Relay error: %v", err)
		}
	}
}

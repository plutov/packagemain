package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxMessage struct {
	ID      uuid.UUID
	Topic   string
	Message []byte
}

func processOutboxMessages(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Lock the next pending message so other relay instances don't grab it
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
	defer rows.Close()

	// 2. If we found a message, publish it to Pub/Sub
	var msg OutboxMessage
	if rows.Next() {
		if err := rows.Scan(&msg.ID, &msg.Topic, &msg.Message); err != nil {
			return err
		}
	} else {
		// No new messages
		return nil
	}

	log.Printf("Publishing message %s to topic %s", msg.ID, msg.Topic)
	// TODO: implement pub/sub message publishing

	// 3. Mark the message as processed
	_, err = tx.Exec(ctx, "UPDATE outbox SET state = 'processed', processed_at = now() WHERE id = $1", msg.ID)
	if err != nil {
		return err
	}
	log.Printf("Marked message %s as processed", msg.ID)

	return tx.Commit(ctx)
}

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if err := processOutboxMessages(context.Background(), pool); err != nil {
			log.Printf("Error processing outbox: %v", err)
		}
	}
}

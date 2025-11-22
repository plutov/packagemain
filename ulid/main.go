package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
)

func main() {
	ctx := context.Background()

	addr := "postgres://postgres:pass@localhost:5432/postgres"
	conn, err := pgx.Connect(ctx, addr)
	if err != nil {
		fmt.Printf("unable to	connect to the database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS ulid_test (
  id UUID PRIMARY KEY,
	kind TEXT NOT NULL,
  value TEXT NOT NULL
);`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to create table: %v\n", err)
		os.Exit(1)
	}

	insertUUID(ctx, conn, "1")
	insertUUID(ctx, conn, "2")
	insertUUID(ctx, conn, "3")
	insertUUID(ctx, conn, "4")
	insertUUID(ctx, conn, "5")

	insertULID(ctx, conn, "1")
	insertULID(ctx, conn, "2")
	insertULID(ctx, conn, "3")
	insertULID(ctx, conn, "4")
	insertULID(ctx, conn, "5")
}

func insertUUID(ctx context.Context, conn *pgx.Conn, value string) error {
	id := uuid.New()
	_, err := conn.Exec(ctx, "INSERT INTO ulid_test (id, value, kind) VALUES ($1, $2, 'uuid')", id, value)
	if err != nil {
		fmt.Printf("unable to insert UUID: %v\n", err)
		return err
	}

	fmt.Printf("Inserted UUID: %s\n", id.String())

	return nil
}

func insertULID(ctx context.Context, conn *pgx.Conn, value string) error {
	id := ulid.Make()
	// as you can see, we don't need to format the ULID as a string, it can be used directly
	_, err := conn.Exec(ctx, "INSERT INTO ulid_test (id, value, kind) VALUES ($1, $2, 'ulid')", id, value)
	if err != nil {
		fmt.Printf("unable to insert ULID: %v\n", err)
		return err
	}

	fmt.Printf("Inserted ULID: %s\n", id.String())

	return nil
}

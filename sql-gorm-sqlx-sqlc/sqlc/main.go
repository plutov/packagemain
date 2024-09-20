package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	conn, err := sql.Open("sqlite3", "./../db.sqlite")
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	users, err := getUsersStats(conn, 2)
	if err != nil {
		log.Fatalf("failed to get users: %v", err)
	}

	fmt.Println(users)
}

func getUsersStats(conn *sql.DB, minPosts int) ([]GetUsersWithMinPostsRow, error) {
	queries := New(conn)

	ctx := context.Background()
	return queries.GetUsersWithMinPosts(ctx)
}

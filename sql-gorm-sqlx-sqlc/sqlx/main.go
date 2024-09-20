package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	conn, err := sqlx.Connect("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	users, err := getUsersStats(conn, 2)
	if err != nil {
		log.Fatalf("failed to get users: %v", err)
	}

	fmt.Println(users)
}

type userStats struct {
	UserName  string `db:"name"`
	PostCount string `db:"post_count"`
}

func getUsersStats(conn *sqlx.DB, minPosts int) ([]userStats, error) {
	users := []userStats{}
	query := `SELECT u.name, COUNT(p.id) AS post_count
FROM users AS u
JOIN posts AS p ON u.id = p.user_id
GROUP BY u.id
HAVING post_count >= ?;`

	if err := conn.Select(&users, query, minPosts); err != nil {
		return nil, err
	}

	return users, nil
}

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	conn, err := sql.Open("sqlite3", "db.sqlite")
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
	UserName  sql.NullString
	PostCount sql.NullInt64
}

func getUsersStats(conn *sql.DB, minPosts int) ([]userStats, error) {
	query := `SELECT u.name, COUNT(p.id) AS post_count
FROM users AS u
JOIN posts AS p ON u.id = p.user_id
GROUP BY u.id
HAVING post_count >= ?;`

	rows, err := conn.Query(query, minPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []userStats{}
	for rows.Next() {
		var user userStats

		if err := rows.Scan(&user.UserName, &user.PostCount); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

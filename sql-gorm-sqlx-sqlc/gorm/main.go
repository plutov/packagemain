package main

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	conn, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	users, err := getUsersStats(conn, 2)
	if err != nil {
		log.Fatalf("failed to get users: %v", err)
	}

	fmt.Println(users)
}

type User struct {
	gorm.Model
	ID    int `gorm:"primaryKey"`
	Name  string
	Posts []Post
}

type Post struct {
	gorm.Model
	ID     int `gorm:"primaryKey"`
	UserID int
}

type userStats struct {
	Name  string
	Count int `gorm:"column:post_count"`
}

func getUsersStats(conn *gorm.DB, minPosts int) ([]userStats, error) {
	var users []userStats
	err := conn.Model(&User{}).
		Select("name", "COUNT(p.id) AS post_count").
		Joins("JOIN posts AS p ON users.id = p.user_id").
		Group("users.id").
		Having("post_count >= ?", minPosts).
		Find(&users).Error
	return users, err
}

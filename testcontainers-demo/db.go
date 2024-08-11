package main

import (
	"database/sql"
	"os"
)

type DB interface {
	Init() error
}

type Postgres struct {
    conn *sql.DB
}

func (p *Postgres) Init() error {
	var err error
	p.conn, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	return err
}

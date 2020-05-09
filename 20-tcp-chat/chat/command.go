package main

const (
	CMD_NICK = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
)

type command struct {
	id     int
	client *client
	args   []string
}

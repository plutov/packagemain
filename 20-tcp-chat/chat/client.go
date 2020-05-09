package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.changeNick(args)
		case "/join":
			c.join(args)
		case "/rooms":
			c.rooms()
		case "/msg":
			c.sendRoomMsg(args)
		case "/quit":
			c.quit()
		default:
			c.err(fmt.Errorf("unknown command: %s", cmd))
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}

func (c *client) changeNick(args []string) {
	c.commands <- command{
		id:     CMD_NICK,
		client: c,
		args:   args,
	}
}

func (c *client) join(args []string) {
	c.commands <- command{
		id:     CMD_JOIN,
		client: c,
		args:   args,
	}
}

func (c *client) sendRoomMsg(args []string) {
	c.commands <- command{
		id:     CMD_MSG,
		client: c,
		args:   args,
	}
}

func (c *client) rooms() {
	c.commands <- command{
		id:     CMD_ROOMS,
		client: c,
	}
}

func (c *client) quit() {
	c.commands <- command{
		id:     CMD_QUIT,
		client: c,
	}
}

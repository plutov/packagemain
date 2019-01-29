package main

import (
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start telnet server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("telnet server started on 0.0.0.0:8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err.Error())
			continue
		}

		go handleConnection(conn)
	}
}

type client struct {
	conn net.Conn
}

func handleConnection(conn net.Conn) {
	log.Printf("client connected: %s", conn.RemoteAddr().String())

	c := &client{
		conn,
	}

	// help message
	c.sendMsgToClient(`
Welcome to TelnetChat!

/nick <name>: get a name, name is mandatory for sending messages
/join <room>: join a room, if room doesn't exist the new room will be created
/say <msg>:   send message to everyone in a room
/quit:        disconnects from the chat server
`)
}

func (c *client) sendMsgToClient(msg string) {
	if _, err := c.conn.Write([]byte(msg + "\n")); err != nil {
		log.Printf("unable to send message to a client: %s", err.Error())
	}
}

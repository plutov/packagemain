package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type client struct {
	conn net.Conn
	name string
	room string
}

type room struct {
	members map[net.Addr]*client
}

var rooms map[string]*room
var roomsMU *sync.RWMutex

const usage = `
/nick <name>: get a name, or stay anonymous
/join <room>: join a room, if room doesn't exist the new room will be created
/say <msg>:   send message to everyone in a room
/quit:        disconnects from the chat server
`

func main() {
	rooms = make(map[string]*room)
	roomsMU = &sync.RWMutex{}

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

func handleConnection(conn net.Conn) {
	log.Printf("client %s connected", conn.RemoteAddr().String())

	c := &client{
		conn: conn,
		name: "anonymous",
	}

	// help message
	c.sendMsgToClient(`
Welcome to TelnetChat!
` + usage)

	c.startChat()
}

func (c *client) sendMsgToClient(msg string) {
	if _, err := c.conn.Write([]byte(msg + "\n")); err != nil {
		log.Printf("unable to send message to a client: %s", err.Error())
	}
}

func (c *client) startChat() {
loop:
	for {
		msg, err := c.readClientInput()
		if err != nil {
			log.Printf("unable to read client input: %s", err.Error())
			continue
		}

		msgArgs := strings.Split(msg, " ")

		switch msgArgs[0] {
		case "/nick":
			if len(msgArgs) < 2 {
				c.sendMsgToClient("usage: /nick <name>")
				break
			}

			name := strings.Join(msgArgs[1:len(msgArgs)], " ")
			c.changeNick(name)
			break
		case "/join":
			if len(msgArgs) < 2 {
				c.sendMsgToClient("usage: /join <room>")
				break
			}

			room := strings.Join(msgArgs[1:len(msgArgs)], " ")
			c.joinRoom(room)
			break
		case "/say":
			if len(msgArgs) < 2 {
				c.sendMsgToClient("usage: /say <msg>")
				break
			}

			if len(c.room) == 0 {
				c.sendMsgToClient("join a room first to send a message")
				return
			}

			message := strings.Join(msgArgs[1:len(msgArgs)], " ")

			roomsMU.Lock()
			defer roomsMU.Unlock()
			rooms[c.room].announce(c, fmt.Sprintf("> %s says: %s", c.name, message))

			break
		case "/quit":
			c.quit()
			break loop
		default:
			c.sendMsgToClient(usage)
		}
	}
}

func (c *client) readClientInput() (string, error) {
	c.conn.Write([]byte(""))
	s, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.Trim(s, "\r\n"), nil
}

func (c *client) changeNick(name string) {
	c.name = name
	c.sendMsgToClient(fmt.Sprintf("all right, I will call you %s", c.name))
}

func (c *client) joinRoom(roomName string) {
	c.quitCurrentRoom()

	c.room = roomName

	roomsMU.Lock()
	defer roomsMU.Unlock()

	_, ok := rooms[c.room]
	// create new room
	if !ok {
		rooms[c.room] = &room{
			members: make(map[net.Addr]*client),
		}
	}

	rooms[c.room].announce(c, fmt.Sprintf("> %s joined the room", c.name))
	rooms[c.room].members[c.conn.RemoteAddr()] = c

	c.sendMsgToClient(fmt.Sprintf("welcome to %s", c.room))
}

func (c *client) quitCurrentRoom() {
	if len(c.room) > 0 {
		roomsMU.Lock()
		defer roomsMU.Unlock()

		rooms[c.room].announce(c, fmt.Sprintf("> %s left the room", c.name))
		delete(rooms[c.room].members, c.conn.RemoteAddr())
	}
}

func (r *room) announce(from *client, msg string) {
	for _, c := range r.members {
		if from.conn.RemoteAddr() != c.conn.RemoteAddr() {
			c.sendMsgToClient(msg)
		}
	}
}

func (c *client) quit() {
	log.Printf("client %s left", c.conn.RemoteAddr())

	c.quitCurrentRoom()

	c.sendMsgToClient("Sad to see you go =(")
	c.conn.Close()
}

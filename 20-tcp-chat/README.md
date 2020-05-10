# packagemain #20: Building a TCP Chat in Go

In this video, we'll be building a TCP chat server using Go, which enables clients to communicate with each other. In this live-coding session, we'll be working with Go's `net` package which very well supports TCP, as well we'll be using channels and goroutines.

Let's refresh our knowledge of what TCP is. TCP(Transmission Control Protocol) is one of the major protocols of the internet, it sits above the network layer and provides a transport mechanism for application layer protocols such as HTTP, SMTP, IRC, etc.

Let's review how our chat will work.

Once the user connects to the chat server using `telnet` command line program, they can use the following commands to talk to the server:

- `/nick <name>` - get a name, otherwise user will stay anonymous.
- `/join <name>` - join a room, if room doesn't exist, the new room will be created. User can be only in one room at the same time.
- `/rooms` - show list of available rooms to join.
- `/msg	<msg>` - broadcast message to everyone in a room.
- `/quit` - disconnects from the chat server.

## Enough talking, show me a code!

The whole application will consist of the following parts:

- client: current user and its connection
- room
- command: from the client to the server
- server: which manages all incoming commands, as well it stores rooms and clients
- TCP server itself to accept network connections

### Main types

Let's start by defining main structures of our chat: client, room, command and server:

#### Command

```go
type commandID int

const (
	CMD_NICK commandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
)

type command struct {
	id     commandID
	client *client
	args   []string
}
```

- `id` - unque command type ID
- `client`- sender of the command
- `args` - slice of strings from client message

#### Client

Client is responsible for keeping user info, TCP connection, as well as parsing user input and sending it to the server via channel.

```go
type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
}
```

- `conn` - client TCP connection
- `nick` - optional nickname, "anonymous" is default value
- `room` - pointer to current room, nil in the beginning
- `commands` - channel of incoming commands, this will be sent to `server` for processing

#### Room

Room holds its name and list of members.

```go
type room struct {
	name    string
	members map[net.Addr]*client
}
```

- `name` - room required name
- `members` - we use client remove address as their unique key, but that may be not the optimal solution

#### Server

Server will be responsible for handling incoming commands, as well for storing the state (rooms at the moment).

```go
type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}
```

- `rooms` - map of rooms
- `commands` - channel for sending commands from the client to the server

### TCP Server

Let's start with building a TCP server, in `main` func we initialize TCP listener and start listening for new messages.

```go
package main

import (
	"log"
	"net"
)

func main() {
    s := newServer()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("server started on :8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}
```

### Client Input

Once new client has connected, we should initialize it and start listening for incoming messages.

```go
func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}

	c.readInput()
}
```

We pass `s.commands` to the client, so later client can send commands to this channel and it will be processed by the server. We can do it, because channels are the "reference" types.

```go
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
			c.commands <- command{
				id:     CMD_NICK,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				id:     CMD_ROOMS,
				client: c,
			}
		case "/msg":
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
			}
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
```

As you noticed `readInput` is a blocking function because it's constantly reading the user input line by line (except when the connection is lost or we have an error). That's why in `main.go` we have a separate goroutine per each client (`go s.newClient(conn)`).

This function doesn't process the input, we'll have `server` do so in centralized manner (also to keep the order of the messages). But it sends a message to the channel for each command received.

### Process messages by server

Our server will have a blocking function `run` which will receive messages and process them:

```go
func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args[1])
		case CMD_JOIN:
			s.join(cmd.client, cmd.args[1])
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		}
	}
}
```

And we will call `run()` function from our `main` func as goroutine.

```go
func main() {
	s := newServer()
    go s.run()

    // ...
}
```

Now what's left is to implement logic for each command:

```go
func (s *server) nick(c *client, nick string) {
	c.nick = nick
	c.msg(fmt.Sprintf("all right, I will call you %s", nick))
}

func (s *server) join(c *client, roomName string) {
	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

	c.msg(fmt.Sprintf("welcome to %s", roomName))
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	msg := strings.Join(args[1:len(args)], " ")
	c.room.broadcast(c, c.nick+": "+msg)
}

func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("sad to see you go =(")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
```

- `nick()` - sets the client's nick and sends confirmation message to the client
- `join()` - joins the room and creates it if it doesn't exist. Note that we don't protect our map with mutex, because all commands are processed synchronously by channel. We also quit current room before joining.
- `listRooms()` - prints current rooms
- `msg()` - broadcasts message to the current room
- `quit()` - closes the connection

```go
func (r *room) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			m.msg(msg)
		}
	}
}
```

Room has `broadcast()` function to send a message to all members of the room.

## Testing

Now it's time to build, run and test it using `telnet` command.

I'll have 3 terminal windows: one for the server, and another 2 for clients.

```
go build .
./chat
server started on :8888
```

Client 1:

```
telnet localhost 8888
/nick john
> all right, I will call you john
/join #general
> welcome to #general
> jack joined the room
> jack: Hi
```

Client 2:

```
telnet localhost 8888
/nick jack
> all right, I will call you jack
/rooms
> available rooms: #general
/join #general
> welcome to #general
/msg Hi
```

## Conclusion

I would like to highlight that this program is not final yet and it misses few very important items. I did this intentionally, so the video doesn't become too long. Some of them are:

- Validation of message body: commands, arguments, body size.
- State: current server is stateless, meaning if it shuts down - all connections will be closed. It can be also accommodated with graceful shutdown.

Feel free to submit a PR with improvements for this TCP chat in [this repo](https://github.com/plutov/packagemain/tree/master/20-tcp-chat/chat).
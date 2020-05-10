# packagemain #19: Bulding a TCP Chat in Go

In this video we'll be building a TCP chat server using Go, which enables clients to communicate with each other. In this live-coding session we'll be working with Go's `net` package which very well supports TCP, as well we'll be using channels and goruotines.

Let's refresh our knowledge of what TCP is. TCP(Transmission Control Protocol) is one of the major protocols of the internet, it sits above the network layer and provides a transport mechanism for application layer protocols such as HTTP, SMTP, IRC,etc.

Let's review how our chat will work.

Once user connects to the chat server using `telnet` command line program, they can use the following commands to talk to the server:

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
- server: which manages all incoming connections and commands, as well it stores rooms and clients
- TCP server itself to accept network connections

### Main types

Let's start by defining main structures of our chat: client, room, command and server:

```go

```

### TCP Server

In `main` func we initialize TCP listener and start listening for new messages.

### Server

Our `server` is the main entity that clients connect to, it will be responsible for storing chat rooms, clients and will be broadcasting messages.

## Conclusion

I would like to highlight that this program is not final yet and it misses few very important items. I did this intentionally, so the video doesn't become too long. Few of them:

- Validation of message body: commands, arguments, body size.
- State: current server doesn't have state, meaning if it shuts down - all connections will be closed. It can be also accomodated with graceful shutdown.
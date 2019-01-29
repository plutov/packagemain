# packagemain #16: Creating simple Telnet Chat

Telnet is one of the earliest remote login protocols on the Internet. It was initally released in the early days of IP networking in 1969, and was for a long time the default way to access remote networked computers.

Telnet clients can connect to Telnet server using `telnet <host> <port>` command, for example type `telnet towel.blinkenlights.nl` and watch Star Wars in ASCII.

In this video I will live code a simple Telnet chat server in Go.

Once user connects to the chat server using Telnet, they can use the following commands to talk to the server:

- `/nick <name>` - reserve a name, name is mandatory for sending messages.
- `/join <room>` - join a room, if room doesn't exist the new room will be created.
- `/say	<msg>` - send message to everyone in a room.
- `/quit` - disconnects from the chat server.
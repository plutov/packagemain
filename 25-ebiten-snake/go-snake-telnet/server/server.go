// Copyright (c) 2017 Alex Pliutau

package server

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/plutov/go-snake-telnet/snake"
)

const (
	leftTopASCII = "\033[0;0H"
	clearASCII   = "\033[2J"
)

// Server struct
type Server struct {
	addr string
}

// New creates new Server instance
func New(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

// Run the telnet server
func (s *Server) Run() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal("Failed to start TCP server: " + err.Error())
	}

	defer listener.Close()
	log.Printf("TCP server started on %s", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %s\n", err.Error())
			continue
		}

		log.Printf("Client: %s\n", conn.RemoteAddr().String())

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	game := snake.NewGame()

	// Clear screen and move to 0:0
	conn.Write([]byte(clearASCII + leftTopASCII))
	conn.Write([]byte(leftTopASCII))

	go s.read(conn, game)
	go game.Start()

	tick := time.Tick(300 * time.Millisecond)
	for range tick {
		// Move to 0:0 and render
		conn.Write([]byte(leftTopASCII + game.Render()))
		if game.IsOver {
			// Cancel ticker
			break
		}
	}

	conn.Close()
}

// Accept input and send it to KeyboardEventsChan
func (s *Server) read(conn net.Conn, game *snake.Game) {
	reader := bufio.NewReader(conn)

	for {
		data, _, err := reader.ReadLine()
		if game.IsOver {
			break
		}
		if err != nil {
			if err == io.EOF {
				game.IsOver = true
				conn.Close()
				break
			}

			log.Println("Read error: " + err.Error())
			continue
		}

		key := strings.ToLower(strings.TrimSpace(string(data)))
		if len(key) > 0 {
			game.KeyboardEventsChan <- snake.KeyboardEvent{
				Key: string(key[0]),
			}
		}
	}
}

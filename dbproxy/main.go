package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const PROXY_ADDR = "127.0.0.1:3307"
const DB_ADDR = "127.0.0.1:3306"
const COM_QUERY = byte(0x03)

func main() {
	proxy, err := net.Listen("tcp", PROXY_ADDR)
	if err != nil {
		fmt.Printf("failed to start proxy: %v", err)
		os.Exit(1)
	}

	for {
		conn, err := proxy.Accept()
		if err != nil {
			fmt.Printf("failed to accept connection: %v", err)
			continue
		}

		fmt.Printf("new connection: %s\n", conn.RemoteAddr())

		go transport(conn)
	}
}

func transport(conn net.Conn) {
	defer conn.Close()

	dbConn, err := net.Dial("tcp", DB_ADDR)
	if err != nil {
		fmt.Printf("failed to connect to db: %v", err)
		return
	}
	defer dbConn.Close()

	// from proxy to db, intercept before forward
	go intercept(conn, dbConn)

	// forward all from db to client, blocking
	if _, err := io.Copy(conn, dbConn); err != nil {
		fmt.Printf("unable to copy: %v", err)
	}
}

func intercept(src, dst net.Conn) {
	// fixed capacity
	buffer := make([]byte, 4096)

	for {
		n, _ := src.Read(buffer)
		if n > 5 {
			// 3 - length of body, 1 - packet sequence number, 1 - command code, etc.
			switch buffer[4] {
			case COM_QUERY:
				query := string(buffer[5:n])
				// todo: use lexer here
				newQuery := strings.ReplaceAll(query, "from orders_v1", "from orders_v2")
				fmt.Printf("original query: %s\n", query)
				fmt.Printf("new query: %s\n", newQuery)

				copy(buffer[5:], []byte(newQuery))
			}

		}

		// forward
		dst.Write(buffer[0:n])
	}
}

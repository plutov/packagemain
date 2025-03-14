package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const PROXY_ADDR = "127.0.0.1:3307"
const DB_ADDR = "127.0.0.1:3306"

func main() {
	proxy, err := net.Listen("tcp", PROXY_ADDR)
	if err != nil {
		log.Fatalf("failed to start proxy: %v", err)
	}

	for {
		conn, err := proxy.Accept()

		log.Printf("new connection: %s", conn.RemoteAddr())
		if err != nil {
			log.Fatalf("failed to accept connection: %v", err)
		}

		go transport(conn)
	}
}

func transport(conn net.Conn) {
	defer conn.Close()

	dbConn, err := net.Dial("tcp", DB_ADDR)
	if err != nil {
		log.Printf("failed to connect to db: %v", err)
		return
	}

	// from proxy to db
	go pipe(dbConn, conn, true)

	// from db to proxy
	pipe(conn, dbConn, false)
}

func pipe(dst net.Conn, src net.Conn, send bool) {
	if send {
		intercept(src, dst)
		return
	}

	_, err := io.Copy(dst, src)
	if err != nil {
		log.Printf("connection error: %s", err.Error())
	}
}

const COM_QUERY = byte(0x03)

func intercept(src, dst net.Conn) {
	buffer := make([]byte, 4096)

	for {
		n, _ := src.Read(buffer)
		if n > 5 {
			switch buffer[4] {
			case COM_QUERY:
				clientQuery := string(buffer[5:n])
				newQuery := rewriteQuery(clientQuery)
				fmt.Printf("client query: %s\n", clientQuery)
				fmt.Printf("server query: %s\n", newQuery)

				writeModifiedPacket(dst, buffer[:5], newQuery)
				continue
			}

		}
		dst.Write(buffer[0:n])
	}
}

func rewriteQuery(query string) string {
	return strings.NewReplacer("from orders_v1", "from orders_v2").Replace(strings.ToLower(query))
}

func writeModifiedPacket(dst net.Conn, header []byte, query string) {
	newBuffer := make([]byte, 5+len(query))
	copy(newBuffer, header)
	copy(newBuffer[5:], []byte(query))
	dst.Write(newBuffer)
}

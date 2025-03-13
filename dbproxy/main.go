package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const PROXY_ADDR = "127.0.0.1:55432"
const PG_ADDR = "127.0.0.1:5432"

func main() {
	proxy, err := net.Listen("tcp", PROXY_ADDR)
	if err != nil {
		log.Fatalf("unable to start proxy: %v", err)
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

func transport(proxyConn net.Conn) {
	defer proxyConn.Close()

	pgConn, err := net.Dial("tcp", PG_ADDR)
	if err != nil {
		log.Printf("failed to connect to db: %v", err)
		return
	}

	fmt.Println("connected to db, proxying...")

	// from proxy to pg
	go pipe(pgConn, proxyConn, true)

	// from pg to proxy
	pipe(proxyConn, pgConn, false)
}

func pipe(dst net.Conn, src net.Conn, send bool) {
	if send {
		intercept(src, dst)
		return
	}

	if _, err := io.Copy(dst, src); err != nil {
		log.Printf("connection error: %s", err.Error())
	}
}

// postgres message format: https://www.postgresql.org/docs/current/protocol-message-formats.html
func intercept(src net.Conn, dst net.Conn) {
	reader := bufio.NewReader(src)
	writer := bufio.NewWriter(dst)

	for {
		messageType, err := reader.ReadByte()
		if err != nil {
			log.Printf("read message type error: %s", err.Error())
			return
		}

		fmt.Printf("message type: %s\n", string(messageType))

		// Read the length of the message (excluding the length itself)
		lengthBytes := make([]byte, 4)
		_, err = io.ReadFull(reader, lengthBytes)
		if err != nil {
			log.Printf("read message length error: %s", err.Error())
			return
		}

		messageLength := int(uint32(lengthBytes[0])<<24 | uint32(lengthBytes[1])<<16 | uint32(lengthBytes[2])<<8 | uint32(lengthBytes[3]))

		// Read the rest of the message
		messageBody := make([]byte, messageLength-4)
		_, err = io.ReadFull(reader, messageBody)
		if err != nil {
			log.Printf("read message error: %s", err.Error())
			return
		}

		if messageType == 'Q' {
			query := string(messageBody)
			newQuery := rewriteQuery(query)

			fmt.Printf("client query: %s\n", query)
			fmt.Printf("server query: %s\n", newQuery)

			sendQuery(writer, newQuery)
		} else {
			writer.WriteByte(messageType)
			writer.Write(lengthBytes)
			writer.Write(messageBody)
		}

		writer.Flush()
	}
}

func rewriteQuery(query string) string {
	return strings.NewReplacer("from orders_v1", "from orders_v2").Replace(strings.ToLower(query))
}

func sendQuery(writer *bufio.Writer, query string) {
	message := []byte(query)
	length := len(message) + 4

	writer.WriteByte('Q')
	writer.Write([]byte{
		byte(length >> 24),
		byte(length >> 16),
		byte(length >> 8),
		byte(length),
	})
	writer.Write(message)
	writer.Flush()
}

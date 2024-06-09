package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	// proxy listens on port 3307
	proxy, err := net.Listen("tcp", ":3307")
	if err != nil {
		log.Fatalf("failed to start proxy: %s", err.Error())
	}

	for {
		conn, err := proxy.Accept()

		log.Printf("new connection: %s", conn.RemoteAddr())
		if err != nil {
			log.Fatalf("failed to accept connection: %s", err.Error())
		}

		go transport(conn)
	}
}

func transport(conn net.Conn) {
	defer conn.Close()

	mysqlAddr := fmt.Sprintf("%s:%s", os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"))
	mysqlConn, err := net.Dial("tcp", mysqlAddr)
	if err != nil {
		log.Printf("failed to connect to mysql: %s", err.Error())
		return
	}

	// from proxy to mysql
	go pipe(mysqlConn, conn, true)
	// from mysql to proxy
	pipe(conn, mysqlConn, false)
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

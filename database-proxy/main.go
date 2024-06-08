package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
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

	readChan := make(chan int64)
	writeChan := make(chan int64)
	var readBytes, writeBytes int64

	// from proxy to mysql
	go pipe(mysqlConn, conn, true)
	// from mysql to proxy
	go pipe(conn, mysqlConn, false)

	readBytes = <-readChan
	writeBytes = <-writeChan

	log.Printf("connection closed. read bytes: %d, write bytes: %d", readBytes, writeBytes)
}

func pipe(dst, src net.Conn, send bool) {
	if send {
		intercept(src, dst)
	}

	_, err := io.Copy(dst, src)
	if err != nil {
		log.Printf("connection error: %s", err.Error())
	}
}

func intercept(src, dst net.Conn) {
	buffer := make([]byte, 4096)

	for {
		n, _ := src.Read(buffer)
		if n > 5 {
			fmt.Println(string(buffer[:n]))
		}
		dst.Write(buffer[0:n])
	}
}

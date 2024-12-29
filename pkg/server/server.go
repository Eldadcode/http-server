package server

import (
	"fmt"
	"log"
	"net"

	"github.com/Eldadcode/http-server/internal/http"
)

const (
	defaultPort uint16 = 8080
)

func readEntirePacket(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func handleConnection(connection net.Conn) {
	defer connection.Close()
	buf, err := readEntirePacket(connection)
	if err != nil {
		_, err = connection.Write([]byte("Error reading packet"))
		if err != nil {
			log.Printf("Failed writing to connection\n")
		}
	}
	httpResponse, err := http.HandleRequest(buf)
	if err != nil {
		_, err = connection.Write([]byte("Error with parsing HTTP request"))
		if err != nil {
			log.Printf("Failed writing to connection\n")
		}
	}

	_, err = connection.Write(httpResponse.Bytes())
	if err != nil {
		log.Printf("Failed writing to connection\n")
	}
}

// StartServer serves an HTTP server on a given port
func StartServer(port uint16) error {
	if port == 0 {
		port = defaultPort
	}
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err.Error())
		return err
	}

	log.Printf("Server started on http://localhost:%d\n", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err.Error())
			return err
		}
		go handleConnection(conn)
	}
}

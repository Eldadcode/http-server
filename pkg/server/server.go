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

func handleConnection(conn net.Conn) {
	var response_buffer []byte = []byte("HTTP/1.0 ")
	defer conn.Close()
	buf, err := readEntirePacket(conn)
	if err != nil {
		conn.Write([]byte("Error reading packet"))
	}
	http_response, err := http.HandleHTTPRequest(buf)
	if err != nil {
		conn.Write([]byte("Error with parsing HTTP request"))
	}

	response_buffer = append(response_buffer, fmt.Sprintf("%s\n", http_response.Status)...)
	response_buffer = append(response_buffer, fmt.Sprintf("Server: %s\n", http_response.Server)...)
	response_buffer = append(response_buffer, fmt.Sprintf("Date: %s\n", http_response.Date)...)
	response_buffer = append(response_buffer, fmt.Sprintf("Content-type: %s\n", http_response.ContentType)...)
	response_buffer = append(response_buffer, fmt.Sprintf("Content-Length: %d\n\n", http_response.ContentLength)...)
	response_buffer = append(response_buffer, http_response.Content...)

	conn.Write(response_buffer)
}

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

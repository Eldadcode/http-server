package server

import (
	"fmt"
	"net"
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
	defer conn.Close()
	buf, err := readEntirePacket(conn)
	if err != nil {
		conn.Write([]byte("Error reading packet"))
	}
	fmt.Printf("Received: %s\n", string(buf))

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

	fmt.Printf("Server started on http://localhost:%d\n", port)
	
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err.Error())
			return err
		}
		go handleConnection(conn)
	}
}

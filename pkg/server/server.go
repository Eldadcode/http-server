package server

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	conn.Read(buf)
	fmt.Printf("Received: %s\n", string(buf))

}

func StartServer(port uint16) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err.Error())
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err.Error())
			return err
		}
		go handleConnection(conn)
	}
}

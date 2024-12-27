package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Eldadcode/http-server/pkg/server"
)

func main() {
	var err error
	var port uint64

	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <port>\n", os.Args[0])
		os.Exit(1)
	}
	port, err = strconv.ParseUint(os.Args[1], 10, 16)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid port number: %s, Error: %s\n", os.Args[1], err.Error())
		os.Exit(1)
	}

	fmt.Printf("Server started on http://localhost:%d\n", port)
	server.StartServer(uint16(port))
}

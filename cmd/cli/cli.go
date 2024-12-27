package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Eldadcode/http-server/pkg/server"
)

func main() {
	var err error
	var port uint64 = 0

	switch {
	case len(os.Args) == 2:
		port, err = strconv.ParseUint(os.Args[1], 10, 16)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid port number: %s, Error: %s\n", os.Args[1], err.Error())
			os.Exit(1)
		}
	case len(os.Args) > 2:
		fmt.Fprintf(os.Stderr, "Usage: %s <Optional(port)>\n", os.Args[0])
		os.Exit(1)
	}

	server.StartServer(uint16(port))
}

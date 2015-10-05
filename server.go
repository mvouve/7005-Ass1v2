package main

import (
	"log"
	"net"
)

const defaultPort = ":7005"

func startServer(port string) {
	if port == "" {
		port = defaultPort
	}

	ln, err := net.Listen("tcp", port)

	if err != nil {
		log.Print(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
		}
		go receiveMessage(conn)
	}
}

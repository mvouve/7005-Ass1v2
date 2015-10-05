package main

import (
	"log"
	"net"
)

func startServer(port string) {
	ln, err := net.Listen("tcp", port)

	if err != nil {
		log.Print(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	for {
		if receiveMessage(conn) != nil {
			return
		}
	}
}

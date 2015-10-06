package main

import (
	"log"
	"net"
)

/******************************************************************************************
 * FUNCTION startServer(port string)
 *
 * DESIGNER: Marc Vouve
 *
 * PROGRAMMER: Marc Vouve
 *
 * DATE: Saturaday October 3rd
 *
 * REVISIONS: None
 *
 *
 * PROTOTYPE: startServer(port string)
 *								port : the port to listen on (NEEDS : in invocation)
 *
 * RETURNS nil
 *
 * NOTES: This function listens for connections and relies on receive message to determine
 * action
 * ***************************************************************************************/
func startServer(port string) {

	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
		}
		go receiveMessage(conn)
	}
}

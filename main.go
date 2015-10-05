package main

import (
	"fmt"
	"log"
	"net"

	"github.com/docopt/docopt-go"
)

func main() {
	usage := `
  Usage:
  stupid_ftp client <host> [(--GET | --SEND)] <FILE>
  stupid_ftp server [<port>]`

	arguments, _ := docopt.Parse(usage, nil, true, "1.0.0", false)

	if arguments["server"].(bool) {
		startServer(arguments["<port>"].(string))
	} else {
		client(arguments["<host>"].(string), arguments["<FILE>"].(string), arguments["--SEND"].(bool))
	}
}

func client(host string, file string, send bool) {
	connect, err := net.Dial("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	if send {
		sendFile(connect, file)
	} else {
		requestFile(connect, file)
		if receiveMessage(connect) != nil {
			fmt.Println("The server did not respond, this likely means that the file does not exist")
		}
	}
}

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/docopt/docopt-go"
)

/******************************************************************************************
 * FUNCTION main
 *
 * DESIGNER: Marc Vouve
 *
 * PROGRAMMER: Marc Vouve
 *
 * DATE: Saturaday October 3rd
 *
 * REVISIONS: Sunday October 4th - simplified protocol
 *
 *
 * PROTOTYPE: main()
 *
 * RETURNS
 *
 * NOTES: This function is the main driver for the program
 * ***************************************************************************************/
func main() {
	usage := `
  Usage:
  stupid_ftp client <host> [(--GET | --SEND)] <FILE>
  stupid_ftp server <port>`

	arguments, _ := docopt.Parse(usage, nil, true, "1.0.0", false)

	if arguments["server"].(bool) {
		startServer(arguments["<port>"].(string))
	} else {
		client(arguments["<host>"].(string), arguments["<FILE>"].(string), arguments["--SEND"].(bool))
	}
}

/******************************************************************************************
 * FUNCTION client(host string, file string, send bool)
 *
 * DESIGNER: Marc Vouve
 *
 * PROGRAMMER: Marc Vouve
 *
 * DATE: Saturaday October 3rd
 *
 * REVISIONS: Sunday October 4th - simplified protocol
 *
 * PROTOTYPE: client(host string, file string, send bool)
 *										host - host and port in x.x.x.x:x format
 *
 * RETURNS
 *
 * NOTES: This function handles the clients flow of control
 * ***************************************************************************************/
func client(host string, file string, send bool) {
	connect, err := net.Dial("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	if send {
		if sendFile(connect, file) != nil {
			log.Fatalln("Could not open file")
		}
	} else {
		requestFile(connect, file)
		if receiveMessage(connect) != nil {
			fmt.Println("The server did not respond, this likely means that the file does not exist")
		}
	}
}

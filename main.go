package main

import (
	"log"

	"github.com/docopt/docopt-go"
)

func main() {
	usage := `
  Usage:
  stupid_ftp client <host> [(--GET | --SEND)] <FILE>
  stupid_ftp server [<port>]`

	arguments, _ := docopt.Parse(usage, nil, true, "0.0.1", false)

	if arguments["server"].(bool) {
		startServer(arguments["<port>"].(string))
	} else {
		client, err := newClient(arguments["<host>"].(string))
		if err != nil {
			log.Fatal(err)
		}
		defer client.close()
		if arguments["--SEND"].(bool) {
			client.send(arguments["<FILE>"].(string))
		} else {
			client.get(arguments["<FILE>"].(string))
		}
	}
}

package main

import "net"

type client struct {
	Conn net.Conn
}

func newClient(host string) (*client, error) {
	connect, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	c := &client{Conn: connect}

	return c, nil
}

func (c *client) close() {
	c.Conn.Close()
}

func (c *client) get(s string) {
	requestFile(c.Conn, s)
	receiveMessage(c.Conn)
}

func (c *client) send(s string) {
	sendFile(c.Conn, s)
}

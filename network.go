package main

import (
	"encoding/binary"
	"io"
	"net"
	"os"
	"strings"
)

const (
	packetSize     = 1024
	packetTypeSize = 2
)

const (
	sendT uint16 = 1
	getT  uint16 = 2
)

func sendFile(conn net.Conn, fn string) {
	buffer := make([]byte, packetSize, packetSize)
	// if the server sent the name of the file, the string can be tainted by 0s that will prevent the file from opening, this splits them off the string
	file, err := os.Open(strings.Trim(fn, string(0)))
	if err != nil {
		conn.Close()
		return
	}
	// send file name
	binary.BigEndian.PutUint16(buffer, sendT)
	copy(buffer[packetTypeSize:], []byte(fn))
	conn.Write(buffer)

	for {
		_, err := file.Read(buffer)
		if err != nil {
			conn.Close()
			return
		}
		conn.Write(buffer)
	}
}

func receiveFile(conn net.Conn, name string) {
	file, _ := os.Create("recv/" + name[0:strings.IndexByte(name, 0)])
	defer file.Close()

	buffer := make([]byte, packetSize, packetSize)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}
		file.Write(buffer[:n])
	}
}

func requestFile(conn net.Conn, fn string) {
	buffer := make([]byte, packetSize, packetSize)
	binary.BigEndian.PutUint16(buffer[0:], getT)
	copy(buffer[packetTypeSize:], []byte(fn))
	conn.Write(buffer)
}

func receiveMessage(conn net.Conn) error {
	charBuf := make([]byte, packetSize, packetSize)
	_, err := conn.Read(charBuf)

	if err == io.EOF {
		return err
	}
	switch binary.BigEndian.Uint16(charBuf) {
	case sendT:
		receiveFile(conn, string(charBuf[packetTypeSize:]))
		return nil
	case getT:
		sendFile(conn, string(charBuf[packetTypeSize:]))
		return nil
	}
	return nil
}

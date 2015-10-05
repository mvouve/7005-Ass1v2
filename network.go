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

/******************************************************************************************
 * FUNCTION sendFile(conn net.Conn, fn string)
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
 * PROTOTYPE: sendFile(conn net.Conn, fn string)
 *								conn - the connection to send the file on
 * 								fn 	 - the name of the file to send
 *
 * RETURNS err - if file does not exist returns error
 *
 * NOTES: this function will close the connection, signalling to the recv end that the
 * 				transmission has ended
 * ***************************************************************************************/
func sendFile(conn net.Conn, fn string) error {
	buffer := make([]byte, packetSize, packetSize)
	// if the server sent the name of the file, the string can be tainted by 0s that will prevent the file from opening, this splits them off the string
	file, err := os.Open(strings.Trim(fn, string(0)))
	if err != nil {
		conn.Close()
		return err
	}
	// send file name
	binary.BigEndian.PutUint16(buffer, sendT)
	copy(buffer[packetTypeSize:], []byte(fn))
	conn.Write(buffer)

	for {
		_, err := file.Read(buffer)
		if err != nil {
			conn.Close()
			return nil
		}
		conn.Write(buffer)
	}
}

/******************************************************************************************
 * FUNCTION receiveFile(conn net.Conn, name string)
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
 * PROTOTYPE: receiveFile(conn net.Conn, fn string)
 *								conn - the connection to send the file on
 * 								fn 	 - the name of the file to send
 *
 * RETURNS nil
 *
 * NOTES: This function is broken only if an error is found on the socket. By sending
 * io.EOF (i.e. close the sender) it will break out of reading
 * ***************************************************************************************/
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

/******************************************************************************************
 * FUNCTION requestFile(conn net.Conn, fn string)
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
 * PROTOTYPE: receiveFile(conn net.Conn, fn string)
 *								conn - the connection to send the file on
 * 								fn 	 - the name of the file to send
 *
 * RETURNS nil
 *
 * NOTES: This function is used by the client to request the file, the server has no need
 * to use this as it doesn't make requests
 * ***************************************************************************************/
func requestFile(conn net.Conn, fn string) {
	buffer := make([]byte, packetSize, packetSize)
	binary.BigEndian.PutUint16(buffer[0:], getT)
	copy(buffer[packetTypeSize:], []byte(fn))
	conn.Write(buffer)
}

/******************************************************************************************
 * FUNCTION receiveMessage(conn net.Conn)
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
 * PROTOTYPE: receiveMessage(conn net.Conn)
 *								conn - the connection to listen on
 *
 * RETURNS nil
 *
 * NOTES: This function reads if the program is reading (sentT) or reading (getT)
 * ***************************************************************************************/
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

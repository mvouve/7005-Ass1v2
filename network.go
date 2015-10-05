package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const packetSize = 1024
const packetBuffer = 1022
const packetTypeSize = 2
const fileSizeSize = 8
const metaHeaderSize = 10

const (
	sendT    uint16 = 1
	requestT uint16 = 2
	metaT    uint16 = 3
)

type packet struct {
	PType   uint16
	Message [packetBuffer]byte
}

type fileInfo struct {
	Size int64
	Name []byte
}

func sendFileInfo(conn net.Conn, file *os.File) {
	buf := make([]byte, packetSize, packetSize)
	finfo, err := file.Stat()
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Print("file not found")
		}
		fmt.Print(err)
	}
	offset := 0
	binary.BigEndian.PutUint16(buf[offset:], metaT)
	offset += 2
	binary.BigEndian.PutUint64(buf[offset:], uint64(finfo.Size()))
	offset += 8
	copy(buf[offset:], finfo.Name())

	conn.Write(buf)
}

func sendFile(conn net.Conn, fn string) {
	buffer := make([]byte, packetSize, packetSize)
	file, err := os.Open(fn[0:strings.Index(fn, string(0))])
	if err != nil {
		return
	}
	sendFileInfo(conn, file)
	for {
		_, err := file.Read(buffer)
		if err != nil {
			return
		}
		conn.Write(buffer)
	}
}

func receiveMessage(conn net.Conn) error {
	charBuf := make([]byte, packetSize, packetSize)
	_, err := conn.Read(charBuf)

	if err == io.EOF {
		return err
	}
	pType := binary.BigEndian.Uint16(charBuf)
	switch pType {
	case sendT:
		break
	case requestT:
		sendFile(conn, string(charBuf[packetTypeSize:]))
		return nil
	case metaT:
		receiveFile(conn, int64(binary.BigEndian.Uint64(charBuf[packetTypeSize:])), string(charBuf[10:]))
		return nil
	}
	return nil
}

func receiveFile(conn net.Conn, size int64, name string) {
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
	binary.BigEndian.PutUint16(buffer[0:], requestT)
	copy(buffer[packetTypeSize:], []byte(fn))

	conn.Write(buffer)
}

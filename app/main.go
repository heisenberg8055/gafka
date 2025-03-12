package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to port 9092")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err.Error())
		os.Exit(1)
	}
	var request []byte
	n, err := conn.Read(request)
	if err != nil {
		log.Fatalf("%v", err.Error())
	}
	fmt.Printf("Read %d bytes: %v\n, %d", n, string(request), len(request))

	messageSize := make([]byte, 4)
	binary.BigEndian.PutUint32(messageSize, 0)

	_, err = conn.Write(messageSize)
	if err != nil {
		log.Fatalf("1:%v", err.Error())
	}
	_, err = conn.Write(request[8:12])
	if err != nil {
		log.Fatalf("2:%v", err.Error())
	}
}

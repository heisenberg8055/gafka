package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	brr := []byte{}
	conn.Read(brr)
	arr := []byte{}
	val := binary.BigEndian.AppendUint32(arr, 7)
	conn.Write(val)
}

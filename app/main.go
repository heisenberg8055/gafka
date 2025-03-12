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
	_, err = conn.Read(brr)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}
	arr := []byte{}
	arr = binary.BigEndian.AppendUint32(arr, 0)
	arr = binary.BigEndian.AppendUint32(arr, 7)
	conn.Write(arr)
}

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
	request := make([]byte, 1024)
	n, err := conn.Read(request)
	if err != nil {
		log.Fatalf("n%v", err.Error())
	}
	fmt.Printf("Read %d bytes: \n", n)

	messageSize := make([]byte, 4)
	binary.BigEndian.PutUint32(messageSize, 0)

	response := make([]byte, 19)

	// // Message_size
	// _, err = conn.Write(request[0:4])
	// if err != nil {
	// 	log.Fatalf("1:%v", err.Error())
	// }

	// Correlation ID
	response = append(response, request[8:12]...)

	// error_code
	binary.BigEndian.PutUint16(response, 0)

	// Tagged Field Encodings
	response = append(response, 2)

	// api keys
	binary.BigEndian.PutUint16(response, 18)

	binary.BigEndian.PutUint16(response, 3)

	binary.BigEndian.PutUint16(response, 4)

	response = append(response, 0)
	binary.BigEndian.PutUint32(response, 0)
	response = append(response, 2)

	_, err = conn.Write(response)
	if err != nil {
		log.Fatalf("1:%v", err.Error())
	}

}

// ApiVersions Response (Version: 4) => error_code [api_keys] throttle_time_ms TAG_BUFFER
//   error_code => INT16
//   api_keys => api_key min_version max_version TAG_BUFFER
//     api_key => INT16
//     min_version => INT16
//     max_version => INT16
//   throttle_time_ms => INT32

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
	var conn net.Conn
	for {
		conn, err = l.Accept()
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

		response := []byte{}
		apiVersion := request[6:8]

		// Message_size
		// _, err = conn.Write(request[0:4])
		// if err != nil {
		// 	log.Fatalf("1:%v", err.Error())
		// }

		apiVer := binary.BigEndian.Uint16(apiVersion)

		switch {
		case apiVer > 0 && apiVer < 5:

			// Correlation ID
			response = append(response, request[8:12]...)

			// error_code
			errorCode := make([]byte, 2)
			binary.BigEndian.PutUint16(errorCode, 0)
			response = append(response, errorCode...)

			response = append(response, 2)

			// api keys
			apiKeyIndex := make([]byte, 2)
			binary.BigEndian.PutUint16(apiKeyIndex, 18)
			response = append(response, apiKeyIndex...)

			apiMin := make([]byte, 2)
			binary.BigEndian.PutUint16(apiMin, 3)
			response = append(response, apiMin...)

			apiMax := make([]byte, 2)
			binary.BigEndian.PutUint16(apiMax, 4)
			response = append(response, apiMax...)

			response = append(response, 0)

			//trottle time
			tt := make([]byte, 4)
			binary.BigEndian.PutUint32(tt, 0)
			response = append(response, tt...)

			response = append(response, 0)

			responseMessageSize := len(response)

			responseMessageArray := make([]byte, 4)
			binary.BigEndian.PutUint32(responseMessageArray, uint32(responseMessageSize))

			_, err = conn.Write(responseMessageArray)
			if err != nil {
				log.Fatalf("1:%v", err.Error())
			}

			_, err = conn.Write(response)
			if err != nil {
				log.Fatalf("1:%v", err.Error())
			}
		default:
			response = append(response, request[8:12]...)
			errorCode := make([]byte, 2)
			binary.BigEndian.PutUint16(errorCode, 35)
			response = append(response, errorCode...)
			responseMessageSize := len(response)

			responseMessageArray := make([]byte, 4)
			binary.BigEndian.PutUint32(responseMessageArray, uint32(responseMessageSize))

			_, err = conn.Write(responseMessageArray)
			if err != nil {
				log.Fatalf("1:%v", err.Error())
			}
			_, err = conn.Write(response)
			if err != nil {
				log.Fatalf("1:%v", err.Error())
			}
		}
	}

}

// ApiVersions Response (Version: 4) => error_code [api_keys] throttle_time_ms TAG_BUFFER
//   error_code => INT16
//   api_keys => api_key min_version max_version TAG_BUFFER
//     api_key => INT16
//     min_version => INT16
//     max_version => INT16
//   throttle_time_ms => INT32

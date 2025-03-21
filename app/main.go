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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	for {
		request := make([]byte, 1024)
		n, err := conn.Read(request)
		if err != nil {
			log.Fatalf("n%v", err.Error())
		}
		fmt.Printf("Read %d bytes: \n", n)

		response := []byte{}
		apiKey := request[4:6]
		apiVersion := request[6:8]

		apiVer := binary.BigEndian.Uint16(apiVersion)
		apiKeyVer := binary.BigEndian.Uint16(apiKey)

		switch apiKeyVer {
		case 18:
			switch {
			case apiVer > 0 && apiVer < 5:

				// Correlation ID
				response = append(response, request[8:12]...)

				// error_code
				errorCode := make([]byte, 2)
				binary.BigEndian.PutUint16(errorCode, 0)
				response = append(response, errorCode...)

				response = append(response, 3)

				// api keys
				apiKeyIndex := make([]byte, 2)
				binary.BigEndian.PutUint16(apiKeyIndex, 18)
				response = append(response, apiKeyIndex...)

				apiMin := make([]byte, 2)
				binary.BigEndian.PutUint16(apiMin, 0)
				response = append(response, apiMin...)

				apiMax := make([]byte, 2)
				binary.BigEndian.PutUint16(apiMax, 4)
				response = append(response, apiMax...)

				response = append(response, 0)

				// DescribeTopicPartitions

				apiKeyIndexD := make([]byte, 2)
				binary.BigEndian.PutUint16(apiKeyIndexD, 75)
				response = append(response, apiKeyIndexD...)

				apiMinD := make([]byte, 2)
				binary.BigEndian.PutUint16(apiMinD, 0)
				response = append(response, apiMinD...)

				apiMaxD := make([]byte, 2)
				binary.BigEndian.PutUint16(apiMaxD, 0)
				response = append(response, apiMaxD...)

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
		case 75:

			response = append(response, request[8:12]...)

			response = append(response, 0)

			tt := make([]byte, 4)
			binary.BigEndian.PutUint32(tt, 0)
			response = append(response, tt...)

			response = append(response, 2)

			// error_code
			errorCode := make([]byte, 2)
			binary.BigEndian.PutUint16(errorCode, 3)
			response = append(response, errorCode...)

			response = append(response, 4)

			topicLength := int(request[25])

			response = append(response, request[26:26+topicLength]...)

			topicID := make([]byte, 16)
			binary.BigEndian.PutUint16(topicID, 0)
			response = append(response, topicID...)

			// isInternal
			response = append(response, 0)

			// partitionArray
			response = append(response, 1)

			// auth op

			auth := make([]byte, 4)
			binary.BigEndian.PutUint16(auth, 0x00000df8)
			response = append(response, auth...)

			response = append(response, 0)

			response = append(response, 0xff)

			response = append(response, 0)

			responseLength := len(response)
			responseMessageArray := make([]byte, 4)
			binary.BigEndian.PutUint32(responseMessageArray, uint32(responseLength))
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

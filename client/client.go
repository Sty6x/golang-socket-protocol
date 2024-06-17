package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net"
	// "time"
)

type Header struct {
	ConnectionType  string
	Namespace       string
	DateEstablished string
	UserId          string
}
type ServerResponseHeader struct {
	ConnectionType  string
	Namespace       string
	ConnectionId    string
	DateEstablished string
	Status          string
}

func main() {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
	}
	socket_header := Header{ConnectionType: "socket", Namespace: "NGEE",
		DateEstablished: "14051239084", UserId: uuid.NewString()}
	encoded_header := encode_request_header(socket_header)
	conn.Write(encoded_header)

	for {
		buffer := make([]byte, 1024)
		bytes_read, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Unable to read server message.")
		}
		decodeResponseHeader := ServerResponseHeader{}
		decode_err := json.Unmarshal(buffer[:bytes_read], &decodeResponseHeader)
		if decode_err != nil {
			fmt.Println("unable to decode")
		}
		fmt.Println(decodeResponseHeader)
	}

}

func encode_request_header(h Header) []byte {
	encoded_header, err := json.Marshal(h)
	if err != nil {
		fmt.Println("Unable to encode socket header")
	}
	decoded_header := Header{}
	decode_err := json.Unmarshal(encoded_header, &decoded_header)
	if decode_err != nil {
		fmt.Println("unable to decode")
	}
	fmt.Printf("Sent: %+v ", decoded_header)

	return encoded_header
}

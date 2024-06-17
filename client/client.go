package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net"
)

type Header struct {
	Protocol        string
	ConnectionType  string
	Namespace       string
	DateEstablished string
	UserId          string
	Payload         Payload
}
type Payload struct{ Data any }

type ServerResponseHeader struct {
	ConnectionType  string
	Namespace       string
	ConnectionId    string
	DateEstablished string
	Status          string
	Payload         Payload
}

func main() {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
	}
	socket_header := Header{ConnectionType: "push", Protocol: "socket", Namespace: "NGEE",
		DateEstablished: "14051239084", UserId: uuid.NewString()}
	encoded_header := encode_request_header(socket_header)
	conn.Write(encoded_header)

	for {
		buffer := make([]byte, 1024)
		bytes_read, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Unable to read server message.")
		}
		serverResponse := ServerResponseHeader{}
		decode_err := json.Unmarshal(buffer[:bytes_read], &serverResponse)
		if decode_err != nil {
			fmt.Println("unable to decode")
		}
		fmt.Println(serverResponse)
		if serverResponse.Status != "OK" {
			fmt.Println("Connection Status: Failed")
			continue
		}
		if serverResponse.ConnectionType == "OK" {
			fmt.Println(serverResponse.ConnectionId)
			fmt.Println(serverResponse.Payload.Data)
		}
	}
}

func encode_request_header(h Header) []byte {
	encoded_header, err := json.Marshal(h)
	if err != nil {
		fmt.Println("Unable to encode socket header")
	}
	return encoded_header
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net"
)

type RequestHeader struct {
	Protocol        string
	ConnectionType  string
	Namespace       string
	DateEstablished string
	UserId          string
}
type Payload struct{ Data any }

type ServerResponseHeader struct {
	ConnectionType  string
	Namespace       string
	ConnectionId    string
	UserId          string
	DateEstablished string
	Status          string
}

type PushHeader struct {
	ConnectionType  string
	Namespace       string
	UserId          string
	DateEstablished string
	Status          string
	Payload         Payload
}

// ConnectionType, Namespace, Payload, ConnectionId

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
	}
	isConnected := establishConnection(conn)

	fmt.Println(isConnected)
	if isConnected {
		fmt.Println("Client connected")
	}
}

func establishConnection(conn net.Conn) bool {
	socket_header := RequestHeader{ConnectionType: "connect", Protocol: "websocket",
		Namespace: "neovim-enjoyers", DateEstablished: "14051239084", UserId: uuid.NewString()}
	encoded_header := encode_request_header(socket_header)
	conn.Write(encoded_header)
	// use this loop to only listen to the server's response of client's initial request

	for {
		buffer := make([]byte, 1024)
		bytes_read, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Unable to read server message.")
			return false
		}
		serverResponse := ServerResponseHeader{}
		decode_err := json.Unmarshal(buffer[:bytes_read], &serverResponse)
		if decode_err != nil {
			fmt.Println("unable to decode")
			continue
		}
		if serverResponse.Status != "OK" {
			fmt.Println("Connection Status: Failed")
			return establishConnection(conn) // reconnect
		}
		return true
	}
}

// use this to get push message from server
// if serverResponse.ConnectionType == "push" {
// 	fmt.Printf("\nServer Message: %s has connected to %s Namespace",
// 		serverResponse.UserId, serverResponse.Namespace)
// }

func encode_request_header(h RequestHeader) []byte {
	encoded_header, err := json.Marshal(h)
	if err != nil {
		fmt.Println("Unable to encode socket header")
	}
	return encoded_header
}

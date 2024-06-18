package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net"
)

type Header struct {
	Protocol       string
	ConnectionType string
}

type Request struct {
	Header
	Namespace       string
	DateEstablished string
	UserId          string
}

type Response struct {
	Header
	ConnectionId    string
	DateEstablished string
	Status          string
}

type PushMessage struct {
	Header
	Status          string
	UserId          string
	Namespace       string
	DateEstablished string
	Payload         string
}

func main() {
	var client net.Conn = initializeClient()
	if client == nil {
		fmt.Println("Unable to connect to the server at this moment.")
	}
	app(client)
}

func app(client net.Conn) {
	for {
		buffer := make([]byte, 1024)
		bytes_read, readErr := client.Read(buffer)
		if readErr != nil {
			fmt.Println("Error occured while reading buffer in the app function")
			break
		}
		pushMessage := PushMessage{}
		pushErr := json.Unmarshal(buffer[:bytes_read], &pushMessage)
		if pushErr != nil {
			fmt.Println("Error occured while decoding push header in the app function")
			break
		}
		if pushMessage.Header.ConnectionType == "push" {
			fmt.Printf("Server Message: User %s has connected in the %s namespace",
				pushMessage.UserId, pushMessage.Namespace)
		}
	}
}

func initializeClient() net.Conn {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
	}
	isConnected := establishConnection(conn)
	if !isConnected {
		fmt.Println("Unable to connect to the server at this moment.")
		return nil
	}
	fmt.Println("\nClient connected")
	return conn
}

func establishConnection(conn net.Conn) bool {
	socket_header := Request{
		Header:    Header{ConnectionType: "connect", Protocol: "websocket"},
		Namespace: "neovim-enjoyers", DateEstablished: "14051239084", UserId: uuid.NewString()}
	encodedHeader := encodeRequestHeader(socket_header)
	conn.Write(encodedHeader)
	for {
		buffer := make([]byte, 1024)
		bytes_read, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Unable to read server message.")
			return false
		}
		serverResponse := Response{}
		decode_err := json.Unmarshal(buffer[:bytes_read], &serverResponse)
		if decode_err != nil {
			fmt.Println("unable to decode")
			continue
		}
		if serverResponse.Header.ConnectionType != "connect" {
			continue
		}
		if serverResponse.Status != "OK" {
			fmt.Println("Connection Status: Failed")
			return establishConnection(conn) // reconnect
		}
		return true
	}
}

func encodeRequestHeader(h Request) []byte {
	encodedHeader, err := json.Marshal(h)
	if err != nil {
		fmt.Println("Unable to encode socket header")
	}
	return encodedHeader
}

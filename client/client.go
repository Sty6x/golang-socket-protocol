package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go-tcp/internal/utils"
	"go-tcp/internal/utils/messageTypes"
	"net"
)

type Header struct {
	Protocol       string
	ConnectionType string
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
		pushMessage := message.PushMessage{}
		pushErr := json.Unmarshal(buffer[:bytes_read], &pushMessage)
		if pushErr != nil {
			fmt.Println("Error occured while decoding push header in the app function")
			break
		}

		// TODO
		// figure out if server push or client push header

		if pushMessage.Header.ConnectionType == "push" {
			fmt.Printf("Server Message: User %s has connected in the %s namespace\n",
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
	socket_header := message.Request{
		Header:    message.Header{ConnectionType: "connect", Protocol: "websocket"},
		Namespace: "neovim-enjoyers", DateEstablished: "14051239084",
		UserId: uuid.NewString()}
	json := &utils.Json{}
	encodedHeader := json.Encode(socket_header)
	conn.Write(encodedHeader)
	for {
		buffer := make([]byte, 1024)
		bytes_read, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Unable to read server message.")
			return false
		}
		serverResponse := json.Decode(buffer[:bytes_read])
		if serverResponse == nil {
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

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go-tcp/internal/utils"
	"go-tcp/internal/utils/message_types"
	"net"
	"os"
)

type User struct {
	conn         net.Conn
	namespace    string
	connectionId string
	userId       string
}

type UserMethods interface {
	PushMessage(inputChan chan string)
}

func (u *User) PushMessage(inputChan chan string) {
	for input := range inputChan {
		if input == "\n" {
			fmt.Printf("\nDont send an empty string %q\n", input)
			continue
		}
		fmt.Printf("\nInput: %q\n", input)
		clientMsg := message.PushMessage{
			Header:          message.Header{Protocol: "websocket", ConnectionType: "push"},
			ConnectionId:    u.connectionId,
			Payload:         input,
			Namespace:       u.namespace,
			DateEstablished: "90123789035478",
		}
		fmt.Printf("Push Message: %+v", clientMsg)
		encodedHeader, jsonErr := json.Marshal(clientMsg)
		if jsonErr != nil {
			fmt.Println("Unable to encode push message")
		}
		// TODO unable to write back to the server, the line below
		// is not able to write for some reason, and does not throw any errors

		// running process is not triggering Write even if we're connected with the server
		_, writeErr := u.conn.Write(encodedHeader)
		if writeErr != nil {
			fmt.Println("Unable to write push message")
		}
	}
}

func main() {
	user := initializeClient()
	if user == nil {
		fmt.Println("Unable to connect to the server at this moment.")
	}
	inputChan := make(chan string)
	go notifyConnection(user.conn) // listening to the server
	go inputLoop(inputChan)
	go user.PushMessage(inputChan)
	fmt.Println("Not connected")
	for {
	}
	// close(inputChan)
}

func inputLoop(inputChan chan string) {
	fmt.Println("Input loop called")
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Unable to read stdin")
			return
		}
		inputChan <- input
	}
}

func notifyConnection(clientConn net.Conn) {
	buffer := make([]byte, 1024)
	bytes_read, readErr := clientConn.Read(buffer)
	if readErr != nil {
		fmt.Println("Error occured while reading buffer in the app function")
	}
	pushMessage := message.PushMessage{}
	pushErr := json.Unmarshal(buffer[:bytes_read], &pushMessage)
	if pushErr != nil {
		fmt.Println("Error occured while decoding push header in the app function")
	}
	if pushMessage.Header.ConnectionType == "push" {
		fmt.Printf("Server Message: User %s has connected in the %s namespace\n",
			pushMessage.UserId, pushMessage.Namespace)
	}
}

func initializeClient() *User {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
	}
	var namespace string = os.Args[1]
	var newId string = uuid.NewString()
	isConnected, connectionId := establishWebsocketConnection(conn,
		message.Request{
			Namespace: namespace,
			UserId:    newId,
		})
	user := User{userId: newId, connectionId: connectionId, namespace: namespace, conn: conn}
	if !isConnected {
		fmt.Println("Unable to connect to the server at this moment.")
		return nil
	}
	fmt.Println("\nClient connected")
	return &user
}

func establishWebsocketConnection(conn net.Conn, msg message.Request) (bool, string) {
	connectMessage := message.Request{
		Header: message.Header{
			ConnectionType: "connect",
			Protocol:       "websocket",
		},
		Namespace:       msg.Namespace,
		DateEstablished: "14051239084",
		UserId:          msg.UserId}
	json := &utils.Json{}
	encodedHeader := json.Encode(connectMessage)
	conn.Write(encodedHeader)
	for {
		buffer := make([]byte, 1024)
		bytes_read, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Unable to read server message.")
			return false, ""
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
			return establishWebsocketConnection(conn,
				message.Request{
					Namespace: msg.Namespace,
					UserId:    msg.UserId,
				}) // reconnect
		}
		return true, serverResponse.ConnectionId
	}
}

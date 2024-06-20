package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go-tcp/internal/users"
	"go-tcp/internal/utils"
	"go-tcp/internal/utils/message_types"
	"net"
	"os"
)

func main() {
	user := initializeClient()
	if user == nil {
		fmt.Println("Unable to connect to the server at this moment.")
	}
	inputChan := make(chan string)
	go inputLoop(inputChan)
	go user.PushMessage(inputChan)
	serverListener(user.Conn) // listening to the server
	close(inputChan)
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

// This thing works it can listen to any connection type from server
func serverListener(clientConn net.Conn) {
	for {
		buffer := make([]byte, 1024)
		bytes_read, readErr := clientConn.Read(buffer)
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

		if pushMessage.Header.ConnectionType == "push" {
			fmt.Printf("\nServer Message: User %s has connected in the %s namespace\n",
				pushMessage.UserId, pushMessage.Namespace)
		}
	}
}

func initializeClient() *users.User {
	conn, err := net.Dial("tcp", "localhost:8080")
	// In  IP6 ::1.<Client Ephemeral Port> > ::1.8080: Flags [S], seq 995990637, win 65476, options [mss 65476,sackOK,TS val 639010819 ecr 0,nop,wscale 7], length 0
	// In  IP6 ::1.8080 > ::1.<Client Ephemeral Port>: Flags [S.], seq 1471586355, ack 995990638, win 65464, options [mss 65476,sackOK,TS val 639010819 ecr 639010819,nop,wscale 7], length 0
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
	fmt.Println(connectionId)
	user := users.User{UserId: newId,
		ConnectionId: connectionId,
		Namespace:    namespace,
		Conn:         conn,
	} // creates the local user on the client side
	if !isConnected {
		fmt.Println("Unable to connect to the server at this moment.")
		return nil
	}
	fmt.Println("\nClient connected")
	return &user
}

// In  IP6 ::1.<Client Ephemeral Port> > ::1.8080: Flags [.], ack 1, win 512, options [nop,nop,TS val 639010819 ecr 639010819], length 0
// In  IP6 ::1.<Client Ephemeral Port> > ::1.8080: Flags [P.], seq 1:151, ack 1, win 512, options [nop,nop,TS val 639010819 ecr 639010819], length 150: HTTP
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
	encodedHeader := json.Encode(connectMessage) // along with ack after tcp handshake
	_, writeErr := conn.Write(encodedHeader)
	if writeErr != nil {
		fmt.Println("Unable to write buffer to the server.")
		return false, "Client error"
	}
	for {
		buffer := make([]byte, 1024)
		bytes_read, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Unable to read server message.")
			return false, ""
		}
		serverResponse := json.Decode(buffer[:bytes_read])
		fmt.Printf("\n%s", serverResponse.Header.ConnectionType)
		if serverResponse == nil {
			fmt.Println("unable to decode")
			break
		}
		// only accept connect connectionType for establishing a connection
		if serverResponse.Header.ConnectionType != "connect" {
			fmt.Println("\nWrong Connection type")
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
	return false, "Disconnect"
}

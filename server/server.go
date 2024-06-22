package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go-tcp/internal/namespaces"
	"go-tcp/internal/users"
	"go-tcp/internal/utils/buffer_utils"
	"go-tcp/internal/utils/message_types"
	"go-tcp/server/websocket"
	"net"
)

const PORT = ":8080"
const BUFFERSIZE = 50

var Namespaces = namespaces.New()
var Users = users.New()

func main() {
	var app = setServerSocket()
	clientMessageBuffer := make(chan message.PushMessage)
	go websocket.RelayClientMessages(clientMessageBuffer)
	CreateTcpConnections(app, clientMessageBuffer)
}

func setServerSocket() net.Listener {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Unable to create a listener ")
	}
	return listener
}

func CreateTcpConnections(server net.Listener, clientMessageBuffer chan message.PushMessage) {
	fmt.Println("Server starts at [::]:8080")
	for {
		conn, err := server.Accept() // Blocks all the process until a new TCP CONNECTION IS ESTABLISHED
		if err != nil {
			fmt.Println("Unable to create a tcp connection")
		}
		buffer := buffer.Decoder(conn)
		if buffer == nil {
			fmt.Println("Unable to read client request buffer.")
			continue
		}
		userTcp, connectionType := establishTcpConnection(conn, buffer)
		if userTcp == nil {
			fmt.Println("Unable to establish tcp connection")
			continue
		}
		if connectionType == "connect" {
			websocket.NewConnectionHandler(userTcp)
		}
		// creates a request listener for every new client connection
		go websocket.RequestListener(conn, clientMessageBuffer)

	}
}

func establishTcpConnection(conn net.Conn, buffer []byte) (*users.User, string) {
	clientRequest := message.Request{}
	err := json.Unmarshal(buffer, &clientRequest)
	if err != nil {
		fmt.Println("Unable to decode client request header")
	}
	fmt.Printf("Log Message: client %s has established a connection\n", clientRequest.UserId)
	user, userExists := Users[clientRequest.UserId]
	if !userExists {
		newUser, connectionType := createUser(clientRequest, conn)
		return &newUser, connectionType
	}
	fmt.Printf("\n%v", Namespaces[user.Namespace])
	return &user, clientRequest.Header.ConnectionType
}

func createUser(clientRequest message.Request, conn net.Conn) (users.User, string) {
	// change this using the namespace directly instaed of looping in the function
	existingUser, userExists := Users[clientRequest.UserId]
	if !userExists {
		newUser := users.User{
			UserId:       clientRequest.UserId,
			Conn:         conn,
			IpAddr:       conn.RemoteAddr().String(),
			Namespace:    clientRequest.Namespace,
			ConnectionId: uuid.NewString(),
		}
		Users[newUser.UserId] = newUser
		ns, nsExist := Namespaces[newUser.Namespace]
		if !nsExist {
			// create new user defined namespace
			Namespaces[newUser.Namespace] = namespaces.Namespace{Name: newUser.Namespace}
		}
		ns = namespaces.Namespace{Name: newUser.Namespace, // append new user's connectionId to the namespace
			ConnectedUsers: append(ns.ConnectedUsers[:],
				newUser.ConnectionId)}
		Namespaces[newUser.Namespace] = ns
		return newUser, clientRequest.ConnectionType
	}
	return existingUser, clientRequest.ConnectionType
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go-tcp/internal/namespaces"
	"go-tcp/internal/users"
	"go-tcp/internal/utils/buffer_utils"
	"go-tcp/internal/utils/message_types"
	"go-tcp/internal/utils/server_utils"
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
	go func(messageBuffer chan message.PushMessage) {
		for clientMessage := range messageBuffer {
			fmt.Printf("\nMessage from %s: %+v\n", clientMessage.UserId, clientMessage.Payload)
			userNamespace, ok := Namespaces[clientMessage.Namespace]
			if !ok {
				fmt.Printf("Namespace does not existc")
				continue
			}
			userNamespace.PushClientMessage(clientMessage)
		}
	}(clientMessageBuffer)
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
		conn, err := server.Accept() // Blocks all the process until there is a TCP CONNECTION IS ESTABLISHED
		fmt.Println("Test")
		if err != nil {
			fmt.Println("Unable to create a tcp connection")
		}
		buffer := buffer.Decoder(conn)
		if buffer == nil {
			fmt.Println("Unable to read client request buffer.")
			continue
		}
		fmt.Println("Message received from client") // does not trigger on
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
	fmt.Printf("Client Request: %+v", clientRequest)
	user, userExists := Users[clientRequest.UserId]
	if !userExists {
		newUser := &users.User{IpAddr: conn.RemoteAddr().String(),
			UserId: clientRequest.UserId, Conn: conn,
			Namespace: clientRequest.Namespace, ConnectionId: uuid.NewString()}
		createUser(newUser)
		return newUser, clientRequest.Header.ConnectionType
	}
	fmt.Printf("\n%v", Namespaces[user.Namespace])
	return &user, clientRequest.Header.ConnectionType
}

func createUser(newUser *users.User) {
	Users[newUser.UserId] = *newUser
	ns, ok := Namespaces[newUser.Namespace]
	if !ok {
		Namespaces[newUser.Namespace] = namespaces.Namespace{Name: newUser.Namespace}
	}

	// change this using the namespace directly instaed of looping in the function
	userExists := server.CheckExistingUserConnnection(ns.ConnectedUsers, newUser.ConnectionId)
	if !userExists {
		ns = namespaces.Namespace{Name: newUser.Namespace, ConnectedUsers: append(ns.ConnectedUsers[:],
			newUser.ConnectionId)}
		Namespaces[newUser.Namespace] = ns
	}
	fmt.Printf("\n%v", Namespaces[newUser.Namespace])
}

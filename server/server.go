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
	go func() {
		for l := range clientMessageBuffer {
			fmt.Printf("\n\n Push message: %+v\n", l)
		}
	}()
	go func() {
		// websocket listener should run indefiniteley
		// to receive messages from clients
		// go websocket.RequestListener(user, connectionType, pushMessageBuffer, messageBuffer)
	}()
	CreateTcpConnections(app)
}

func setServerSocket() net.Listener {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Unable to create a listener ")
	}
	return listener
}

func CreateTcpConnections(server net.Listener) {
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
			websocket.SendWebsocketConnectionID(userTcp)
			ns, _ := Namespaces[userTcp.Namespace]
			go ns.NotifyNamespaceUsers(userTcp)
		}
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

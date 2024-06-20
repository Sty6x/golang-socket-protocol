package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go-tcp/internal/namespaces"
	"go-tcp/internal/utils/buffer_utils"
	"go-tcp/internal/utils/global_types"
	"go-tcp/internal/utils/message_types"
	"go-tcp/internal/utils/server_utils"
	"go-tcp/internal/websocket"
	"net"
)

const PORT = ":8080"
const BUFFERSIZE = 50

var Namespaces = namespaces.New()
var Users = make(map[string]types.User)

func main() {
	var app = setServerSocket()
	clientMessageBuffer := make(chan message.PushMessage, BUFFERSIZE)
	go func() {
		for l := range clientMessageBuffer {
			fmt.Printf("\n\n Push message: %+v\n", l)
			// send the client messages in the channel buffer
		}
	}()
	startServer(app, clientMessageBuffer)
	app.Close()
}

func setServerSocket() net.Listener {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Unable to create a listener ")
	}
	return listener
}

// server should handle listening to every connection types from the client
func startServer(server net.Listener, messageBuffer chan message.PushMessage) {

	fmt.Println("Server start at localhost:8080")
	// Listens, Reads and Writes to the client.
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Unable to create a tcp connection")
		}
		buffer := buffer.Decoder(conn)
		if buffer != nil {
			fmt.Println("Unable to read client request buffer.")
		}
		userTcp, connectionType := establishTcpConnection(conn, buffer)
		if userTcp == nil {
			fmt.Println("Unable to establish tcp connection")
			continue
		}
		websocket.RequestListeners(userTcp, connectionType, buffer, messageBuffer)
	}
}

func successSocketConnection(user *types.User) {
	serverResponseHeader := message.Response{
		Header:          message.Header{Protocol: "Websocket", ConnectionType: "connect"},
		DateEstablished: "412908124",
		Status:          "OK",
		ConnectionId:    user.ConnectionId,
	}
	encodedHeader, err := json.Marshal(serverResponseHeader)
	if err != nil {
		fmt.Println("Unable to encode server response header.")
	}
	user.Conn.Write(encodedHeader)
}

func establishTcpConnection(conn net.Conn, buffer []byte) (*types.User, string) {
	parsedHeader := parseClientRequest(buffer)
	user, userExists := Users[parsedHeader.UserId]
	if !userExists {
		newUser := &types.User{IpAddr: conn.RemoteAddr().String(),
			UserId: parsedHeader.UserId, Conn: conn,
			Namespace: parsedHeader.Namespace, ConnectionId: uuid.NewString()}
		createUser(newUser)
		return newUser, parsedHeader.ConnectionType
	}
	fmt.Printf("\n%v", Namespaces[user.Namespace])
	return &user, parsedHeader.ConnectionType
}

func createUser(newUser *types.User) {
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

func parseClientRequest(h []byte) *message.Request {
	clientRequestHeader := message.Request{}
	err := json.Unmarshal(h, &clientRequestHeader)
	if err != nil {
		fmt.Println("Unable to decode client request header")
	}
	return &clientRequestHeader
}

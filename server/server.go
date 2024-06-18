package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go-tcp/internal/utils/message_types"
	"go-tcp/internal/utils/server_utils"
	"net"
)

type User struct {
	ipAddr       string
	conn         net.Conn
	namespace    string
	connectionId string
	userId       string
}

type Namespace struct {
	name           string
	connectedUsers []string
}

const PORT = ":8080"

var Namespaces = make(map[string]Namespace)
var Users = make(map[string]User)

func main() {
	var app = setServerSocket()
	startServer(app)
}

func (ns *Namespace) notifyUsers(userTcp *User) {
	responseHeader := message.PushMessage{
		Header:    message.Header{Protocol: "Websocket", ConnectionType: "push"},
		Namespace: userTcp.namespace,
		Status:    "OK",
		UserId:    userTcp.userId,
	}
	encodedHeader, err := json.Marshal(responseHeader)
	if err != nil {
		fmt.Println("Unable to encode notification header")
		return
	}
	for _, user := range Users {
		if user.namespace == ns.name && user.userId != userTcp.userId {
			fmt.Printf("\nConnection: %+v \n", user)
			user.conn.Write(encodedHeader)
		}
	}
}

func setServerSocket() net.Listener {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Unable to create a listener ")
	}
	return listener
}

func startServer(server net.Listener) {

	fmt.Println("Server start at localhost:8080")
	// Listens, Reads and Writes to the client.
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Unable to create a tcp connection")
		}
		userTcp, connectionType := handleTcpConnection(conn)
		if userTcp != nil {
			if connectionType == "connect" {
				establishSocketConnection(userTcp)
				ns, _ := Namespaces[userTcp.namespace]
				go ns.notifyUsers(userTcp)
			}
			// TODO
			// handle connected client messages
		}
	}
}

func establishSocketConnection(user *User) {
	serverResponseHeader := message.Response{
		Header:          message.Header{Protocol: "Websocket", ConnectionType: "connect"},
		DateEstablished: "412908124",
		Status:          "OK",
		ConnectionId:    user.connectionId,
	}
	encodedHeader, err := json.Marshal(serverResponseHeader)
	if err != nil {
		fmt.Println("Unable to encode server response header.")
	}
	user.conn.Write(encodedHeader)
}

func handleTcpConnection(conn net.Conn) (*User, string) {
	buffer := make([]byte, 1024)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Cant Reaaaad")
		return nil, ""
	}
	parsedHeader := parseClientRequest(buffer[:bytesRead])
	user, userExists := Users[parsedHeader.UserId]
	if !userExists {
		newUser := &User{ipAddr: conn.RemoteAddr().String(),
			userId: parsedHeader.UserId, conn: conn,
			namespace: parsedHeader.Namespace, connectionId: uuid.NewString()}
		createUser(newUser)
		return newUser, parsedHeader.ConnectionType
	}
	fmt.Printf("\n%v", Namespaces[user.namespace])
	return &user, parsedHeader.ConnectionType
}

func createUser(newUser *User) {
	Users[newUser.userId] = *newUser
	ns, ok := Namespaces[newUser.namespace]
	if !ok {
		Namespaces[newUser.namespace] = Namespace{name: newUser.namespace}
	}

	userExists := server.CheckExistingUserConnnection(ns.connectedUsers, newUser.connectionId)
	if !userExists {
		ns = Namespace{name: newUser.namespace, connectedUsers: append(ns.connectedUsers[:],
			newUser.connectionId)}
		Namespaces[newUser.namespace] = ns

	}
	fmt.Printf("\n%v", Namespaces[newUser.namespace])
}

func parseClientRequest(h []byte) *message.Request {
	clientRequestHeader := message.Request{}
	err := json.Unmarshal(h, &clientRequestHeader)
	if err != nil {
		fmt.Println("Unable to decode client request header")
	}
	return &clientRequestHeader
}

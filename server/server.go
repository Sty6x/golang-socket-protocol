package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go-tcp/internal/utils"
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

type ClientConnectHeader struct {
	ConnectionType  string
	Namespace       string
	DateEstablished string
	UserId          string
}

type ServerResponseHeader struct {
	ConnectionType  string
	Namespace       string
	ConnectionId    string
	DateEstablished string
	Status          string
}

const PORT = ":8080"

var Namespaces = make(map[string]Namespace)
var Users = make(map[string]User)

func main() {
	var app = setServerSocket()
	startServer(app)
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
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Unable to create a tcp connection")
		}
		userTcp := handleTcpConnection(conn)
		fmt.Printf("Recevied: %+v ", userTcp)
		if userTcp != nil {
			establishSocketConnection(userTcp)
		}
	}
}

func establishSocketConnection(user *User) {
	serverResponseHeader := ServerResponseHeader{Namespace: user.namespace,
		DateEstablished: "412908124", Status: "OK", ConnectionId: user.connectionId,
		ConnectionType: "socket"}
	encodedHeader, err := json.Marshal(serverResponseHeader)
	if err != nil {
		fmt.Println("Unable to encode server response header.")
	}
	user.conn.Write(encodedHeader)
}

func handleTcpConnection(conn net.Conn) *User {
	buffer := make([]byte, 1024)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Cant Reaaaad")
		return nil
	}
	parsedHeader := parseJsonHeader(buffer[:bytesRead])
	user, userExists := Users[parsedHeader.UserId]
	if !userExists {
		newUser := &User{ipAddr: conn.RemoteAddr().String(),
			userId: parsedHeader.UserId, conn: conn,
			namespace: parsedHeader.Namespace, connectionId: uuid.NewString()}
		createUser(newUser)
		return newUser
	}
	return &user
}

func createUser(newUser *User) {
	Users[newUser.userId] = *newUser
	ns, ok := Namespaces[newUser.namespace]
	if !ok {
		Namespaces[newUser.namespace] = Namespace{name: newUser.namespace}
	}

	userExists := utils.CheckExistingUserConnnection(ns.connectedUsers, newUser.connectionId)
	if !userExists {
		ns = Namespace{name: newUser.namespace, connectedUsers: append(ns.connectedUsers[:], newUser.connectionId)}
		Namespaces[newUser.namespace] = ns
		fmt.Println(Namespaces[newUser.namespace])
	}
}

func parseJsonHeader(h []byte) *ClientConnectHeader {
	clientRequestHeader := ClientConnectHeader{}
	err := json.Unmarshal(h, &clientRequestHeader)
	if err != nil {
		fmt.Println("Unable to decode client request header")
	}
	return &clientRequestHeader
}

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
const BUFFERSIZE = 50

var Namespaces = make(map[string]Namespace)
var Users = make(map[string]User)

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

func bufferReader(conn net.Conn) []byte {
	buffer := make([]byte, 1024)
	bytesRead, buffErr := conn.Read(buffer)
	if buffErr != nil {
		fmt.Println("Cant Reaaaad")
		return nil
	}
	return buffer[:bytesRead]
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
		buffer := bufferReader(conn)
		if buffer != nil {
			fmt.Println("Unable to read client request buffer.")
		}
		userTcp, connectionType := establishTcpConnection(conn, buffer)
		if userTcp != nil {
			if connectionType == "connect" {
				successSocketConnection(userTcp)
				ns, _ := Namespaces[userTcp.namespace]
				go ns.notifyUsers(userTcp)
			}
			if connectionType == "push" {
				clientMsg := message.PushMessage{}
				err := json.Unmarshal(buffer, &clientMsg)
				if err != nil {
					fmt.Println("Unable to decode client request header")
				}
				messageBuffer <- clientMsg
			}
		}
	}
}

func successSocketConnection(user *User) {
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

func establishTcpConnection(conn net.Conn, buffer []byte) (*User, string) {
	parsedHeader := parseClientRequest(buffer)
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

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

type ConnectHeader struct {
	ConnectionType   string
	Namespace        string
	Date_established string
	UserId           string
}

type ConnectResponseHeader struct {
	Connection_type  string
	Namespace        string
	Date_established string
	User_id          string
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
		handleTcpConnection(conn)
	}

}

func handleTcpConnection(conn net.Conn) {
	buffer := make([]byte, 1024)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Cant Reaaaad")
	}
	parsedHeader := parseJsonHeader(buffer[:bytesRead])
	_, userExists := Users[parsedHeader.UserId]
	if !userExists {
		newUser := &User{ipAddr: conn.RemoteAddr().String(),
			userId: parsedHeader.UserId, conn: conn,
			namespace: parsedHeader.Namespace, connectionId: uuid.NewString()}
		createUser(newUser)
	}

}

func createUser(newUser *User) {
	Users[newUser.userId] = *newUser
	ns, ok := Namespaces[newUser.namespace]
	if !ok {
		Namespaces[newUser.namespace] = Namespace{name: newUser.namespace}
	}

	userExists := utils.CheckExistingUserConnnection(ns.connectedUsers, newUser.connectionId)
	if !userExists {
		ns = Namespace{name: newUser.namespace, connectedUsers: append(ns.connectedUsers[:], newUser.userId)}
		Namespaces[newUser.namespace] = ns
	}
	fmt.Println(Namespaces[newUser.namespace])
}

func CheckExistingUserConnnection(connections []string, target string) bool {
	for _, id := range connections {
		if id == target {
			return true
		}
	}
	return false
}

func parseJsonHeader(h []byte) *ConnectHeader {
	client_request_header := ConnectHeader{}
	err := json.Unmarshal(h, &client_request_header)
	if err != nil {
		fmt.Println("Unable to decode client request header")
	}
	return &client_request_header
}

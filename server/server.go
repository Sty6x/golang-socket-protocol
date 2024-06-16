package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go-tcp/internal/utils"
	"net"
)

type User struct {
	ip_addr       string
	conn          net.Conn
	namespace     string
	connection_id string
	user_id       string
}

type Namespace struct {
	name            string
	connected_users []string
}

type Header struct {
	Connection_type  string
	Namespace        string
	Date_established string
	User_id          string
}

const PORT = ":8080"

var Namespaces = make(map[string]Namespace)
var Users = make(map[string]User)

func main() {
	var app = set_server_socket()
	start_server(app)
}

func set_server_socket() net.Listener {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Unable to create a listener ")
	}
	return listener
}

func start_server(server net.Listener) {

	fmt.Println("Server start at localhost:8080")
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Unable to create a tcp connection")
		}
		handle_tcp_connection(conn)
	}

}

func handle_tcp_connection(conn net.Conn) {
	buffer := make([]byte, 1024)
	bytes_read, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Cant Reaaaad")
	}

	parsed_header := parse_json_header(buffer[:bytes_read])
	new_user := &User{ip_addr: conn.RemoteAddr().String(),
		user_id: parsed_header.User_id, conn: conn,
		namespace: parsed_header.Namespace, connection_id: uuid.NewString()}
	create_user(new_user)
}

func create_user(new_user *User) {
	Users[new_user.user_id] = *new_user
	ns, ok := Namespaces[new_user.namespace]
	if !ok {
		Namespaces[new_user.namespace] = Namespace{name: new_user.namespace}
	}

	user_exists := utils.Check_existing_user_connection(ns.connected_users, new_user.connection_id)
	if !user_exists {
		ns = Namespace{name: new_user.namespace, connected_users: append(ns.connected_users[:], new_user.user_id)}
		Namespaces[new_user.namespace] = ns
	}
	fmt.Println(Namespaces[new_user.namespace])
}

func check_existing_user_connection(connections []string, target string) bool {
	for _, id := range connections {
		if id == target {
			return true
		}
	}
	return false
}

func parse_json_header(h []byte) *Header {
	client_request_header := Header{}
	err := json.Unmarshal(h, &client_request_header)
	if err != nil {
		fmt.Println("Unable to decode client request header")
	}
	return &client_request_header
}

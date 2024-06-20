package types

import (
	"net"
)

type User struct {
	IpAddr       string
	Conn         net.Conn
	Namespace    string
	ConnectionId string
	UserId       string
}

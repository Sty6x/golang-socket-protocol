package users

import (
	"encoding/json"
	"fmt"
	message "go-tcp/internal/utils/message_types"
	"net"
	"sync"
)

type User struct {
	IpAddr       string
	Conn         net.Conn
	Namespace    string
	ConnectionId string
	UserId       string
}

type UserMethods interface {
	PushMessage(inputChan chan string)
}
type UsersContainer map[string]User

var (
	once     sync.Once
	instance UsersContainer
)

func New() UsersContainer {
	once.Do(func() {
		instance = make(UsersContainer)
	})
	return instance
}

func (u *User) PushMessage(inputChan chan string) {
	for input := range inputChan {
		if input == "\n" {
			fmt.Printf("\nDont send an empty string %q\n", input)
			continue
		}
		clientMsg := message.PushMessage{
			Header: message.Header{
				Protocol:       "websocket",
				ConnectionType: "push",
			},
			ConnectionId:    u.ConnectionId,
			UserId:          u.UserId,
			Payload:         input,
			Namespace:       u.Namespace,
			DateEstablished: "90123789035478",
		}
		encodedHeader, jsonErr := json.Marshal(clientMsg)
		if jsonErr != nil {
			fmt.Println("Unable to encode push message")
		}
		_, writeErr := u.Conn.Write(encodedHeader)
		if writeErr != nil {
			fmt.Println("Unable to write push message")
		}
	}
}

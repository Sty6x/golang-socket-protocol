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
		fmt.Printf("\nInput: %q\n", input)
		clientMsg := message.PushMessage{
			Header: message.Header{
				Protocol:       "websocket",
				ConnectionType: "connect",
			},
			ConnectionId:    u.ConnectionId,
			Payload:         input,
			Namespace:       u.Namespace,
			DateEstablished: "90123789035478",
		}
		fmt.Printf("\nPush Message: %+v\n", clientMsg)
		encodedHeader, jsonErr := json.Marshal(clientMsg)
		if jsonErr != nil {
			fmt.Println("Unable to encode push message")
		}
		bytesWritten, writeErr := u.Conn.Write(encodedHeader)
		fmt.Printf("Sent: %d\n", bytesWritten)
		if writeErr != nil {
			fmt.Println("Unable to write push message")
		}
	}
}

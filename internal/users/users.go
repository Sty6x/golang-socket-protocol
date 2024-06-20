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
			Header:          message.Header{Protocol: "websocket", ConnectionType: "push"},
			ConnectionId:    u.ConnectionId,
			Payload:         input,
			Namespace:       u.Namespace,
			DateEstablished: "90123789035478",
		}
		fmt.Printf("Push Message: %+v", clientMsg)
		encodedHeader, jsonErr := json.Marshal(clientMsg)
		if jsonErr != nil {
			fmt.Println("Unable to encode push message")
		}
		// TODO unable to write back to the server, the line below
		// is not able to write for some reason, and does not throw any errors

		// running process is not triggering Write even if we're connected with the server
		_, writeErr := u.Conn.Write(encodedHeader)
		if writeErr != nil {
			fmt.Println("Unable to write push message")
		}
	}
}

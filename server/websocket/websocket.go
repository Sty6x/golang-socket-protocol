package websocket

import (
	"encoding/json"
	"fmt"
	"go-tcp/internal/namespaces"
	"go-tcp/internal/users"
	"go-tcp/internal/utils/message_types"
)

func RequestListeners(userTcp *users.User, connectionType string, buffer []byte,
	messageBuffer chan message.PushMessage) {
	NS := namespaces.New()
	if connectionType == "connect" {
		onSuccess(userTcp)
		ns, _ := NS[userTcp.Namespace]
		go ns.NotifyUsers(userTcp)
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
func onSuccess(user *users.User) {
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

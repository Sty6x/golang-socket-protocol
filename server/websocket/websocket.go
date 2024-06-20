package websocket

import (
	"encoding/json"
	"fmt"
	"go-tcp/internal/namespaces"
	"go-tcp/internal/users"
	"go-tcp/internal/utils/message_types"
)

// This thing works it can listen to any connection type from the client.
func RequestListener(userTcp *users.User, connectionType string, buffer []byte,
	messageBuffer chan message.PushMessage) {
	NS := namespaces.New()
	if connectionType == "push" {
		clientMsg := message.PushMessage{}
		err := json.Unmarshal(buffer, &clientMsg)
		if err != nil {
			fmt.Println("Unable to decode client request header")
		}
		fmt.Println("Received")
		messageBuffer <- clientMsg
		return
	}
	if connectionType == "connect" {
		sendWebsocketConnectionID(userTcp)
		ns, _ := NS[userTcp.Namespace]
		go ns.NotifyNamespaceUsers(userTcp)
		return
	}

}

// sends the a websocket connection ID back to client
// In IP6 ::1.8080 > ::1.<Client Ephemeral Port>: Flags [P.], seq 1:150, ack 151, win 512, options [nop,nop,TS val 639010820 ecr 639010819], length 149: HTTP
func sendWebsocketConnectionID(user *users.User) {
	serverResponseHeader := message.Response{
		Header:          message.Header{Protocol: "websocket", ConnectionType: "connect"},
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

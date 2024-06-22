package websocket

import (
	"encoding/json"
	"fmt"
	"go-tcp/internal/namespaces"
	"go-tcp/internal/users"
	"go-tcp/internal/utils/message_types"
	"net"
)

// This thing works it can listen to any connection type from the client.
// @Params
// - user: to read a user requests from user.Conn.Read().
// - pushMessageBuffer: to pass data through the channel buffer for pushing the messages.
func RequestListener(client net.Conn, messageBuffer chan message.PushMessage) {
	for {
		pushMessageBuffer := make([]byte, 1024)
		bytesRead, err := client.Read(pushMessageBuffer)
		if err != nil {
			fmt.Println("Error reading pushed message from websocket listener.")
			break
		}
		clientMsg := message.PushMessage{}
		jsonErr := json.Unmarshal(pushMessageBuffer[:bytesRead], &clientMsg)
		if jsonErr != nil {
			fmt.Println("Unable to decode client request header")
			break
		}
		if clientMsg.Header.ConnectionType == "push" {
			messageBuffer <- clientMsg
		}
	}
}

// sends the a websocket connection ID back to client
// In IP6 ::1.8080 > ::1.<Client Ephemeral Port>: Flags [P.], seq 1:150, ack 151, win 512, options [nop,nop,TS val 639010820 ecr 639010819], length 149: HTTP
func SendWebsocketConnectionID(user *users.User) {
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

func NewConnectionHandler(user *users.User) {
	Namespaces := namespaces.New()
	SendWebsocketConnectionID(user)
	ns, _ := Namespaces[user.Namespace]
	go ns.NotifyNamespaceUsers(user)
}

func RelayClientMessages(messageBuffer chan message.PushMessage) {
	Namespaces := namespaces.New()
	for clientMessage := range messageBuffer {
		fmt.Printf("\nLog Message: A Push Message from %s to %s the namespace", clientMessage.UserId, clientMessage.Namespace)
		userNamespace, ok := Namespaces[clientMessage.Namespace]
		if !ok {
			fmt.Printf("Namespace does not existc")
			continue
		}
		userNamespace.PushClientMessage(clientMessage)
	}
}

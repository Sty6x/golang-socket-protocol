package namespaces

import (
	"encoding/json"
	"fmt"
	"go-tcp/internal/users"
	"go-tcp/internal/utils/global_types"
	"go-tcp/internal/utils/message_types"
	"sync"
)

type Namespace struct {
	Name           string
	ConnectedUsers []string
}

type NamespaceMethods interface {
	notifyUsers(userTcp *types.User)
}

type NamespaceContainer map[string]Namespace

var (
	instance NamespaceContainer
	once     sync.Once
)

func New() NamespaceContainer {
	once.Do(func() {
		instance = make(NamespaceContainer)
	})
	return instance
}

func (ns *Namespace) NotifyUsers(userTcp *types.User) {
	users := users.New()
	responseHeader := message.PushMessage{
		Header:    message.Header{Protocol: "Websocket", ConnectionType: "push"},
		Namespace: userTcp.Namespace,
		Status:    "OK",
		UserId:    userTcp.UserId,
	}
	encodedHeader, err := json.Marshal(responseHeader)
	if err != nil {
		fmt.Println("Unable to encode notification header")
		return
	}
	for _, user := range users {
		if user.Namespace == ns.Name && user.UserId != userTcp.UserId {
			fmt.Printf("\nConnection: %+v \n", user)
			user.Conn.Write(encodedHeader)
		}
	}
}

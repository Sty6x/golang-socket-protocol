package users

import (
	types "go-tcp/internal/utils/global_types"
	"sync"
)

type UsersContainer map[string]types.User

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

package utils

import (
	"encoding/json"
	"fmt"
	"go-tcp/internal/utils/message_types"
)

type JsonParser interface {
	Encode(r message.Request) []byte
	Decode(r []byte) message.Response
}
type Json struct{}

func (j *Json) Encode(r message.Request) []byte {
	encodedHeader, err := json.Marshal(r)
	if err != nil {
		fmt.Println("Unable to encode socket header")
	}
	return encodedHeader
}

func (j *Json) Decode(r []byte) *message.Response {
	response := message.Response{}
	err := json.Unmarshal(r, &response)
	if err != nil {
		fmt.Println("Unable to encode socket header")
		return nil
	}
	return &response
}

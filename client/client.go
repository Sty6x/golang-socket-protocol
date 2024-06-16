package main

import (
	"encoding/json"
	"fmt"
	"net"
	// "time"
)

type Header struct {
	Connection_type  string
	Namespace        string
	Date_established string
	User_id          string
}

func main() {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
	}
	socket_header := Header{Connection_type: "socket", Namespace: "Books", Date_established: "14051239084", User_id: "1239-0asdlkjsd"}
	encoded_header := encode_request_header(socket_header)
	conn.Write(encoded_header)

}

func encode_request_header(h Header) []byte {
	encoded_header, err := json.Marshal(h)
	if err != nil {
		fmt.Println("Unable to encode socket header")
	}
	decoded_header := Header{}
	decode_err := json.Unmarshal(encoded_header, &decoded_header)
	if decode_err != nil {
		fmt.Println("unable to decode")
	}
	fmt.Println(decoded_header)

	return encoded_header
}

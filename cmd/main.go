package main

import (
	"fmt"
	"net"
)

func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		buffer := make([]byte, 3)
		bytes_read, err := conn.Read(buffer)
		if err != nil {
			return
		}
		fmt.Printf("Bytes read: %d", bytes_read)
		// message := string(buffer[:bytes_read])
		// fmt.Println(message)
	}

}

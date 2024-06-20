package buffer

import (
	"fmt"
	"net"
)

func Decoder(conn net.Conn) []byte {
	buffer := make([]byte, 1024)
	bytesRead, buffErr := conn.Read(buffer)
	if buffErr != nil {
		fmt.Println("Cant Reaaaad")
		return nil
	}
	return buffer[:bytesRead]
}

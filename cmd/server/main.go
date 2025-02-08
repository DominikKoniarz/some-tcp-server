package main

import (
	"fmt"

	server "github.com/DominikKoniarz/some-tcp-server/internal/server"
)

func main() {
	protocolVersion := "0001"
	messageType := "I"
	data := "db1\x00user1"
	length := fmt.Sprintf("%04d", 4+4+1+len(data))

	_ = []byte(length + protocolVersion + messageType + data)

	server := server.NewServer()
	server.Start()

	defer server.Stop()
}

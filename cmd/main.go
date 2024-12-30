package main

import (
	"bytes"
	"fmt"
)

type ClientRequest struct {
	command string
	data    string
}

func ParseRequest(rawData []byte) {
	// 4 bytes for length of whole command
	// 4 bytes for protocol version
	// 4 bytes for message type
	// rest of the bytes for data

	// 4 bytes for length of whole command
	length := rawData[:4]
	// 4 bytes for protocol version
	protocolVersion := rawData[4:8]
	// 4 bytes for message type
	messageType := rawData[8:12]
	// rest of the bytes for data
	data := rawData[12:]

	fmt.Printf("Length: %v\n", string(length))
	fmt.Printf("Protocol Version: %v\n", string(protocolVersion))
	fmt.Printf("Message Type: %v\n", string(messageType))
	fmt.Printf("Data: %v\n", string(data))

	stringMessageType := string(messageType)

	if stringMessageType == "INIT" {
		// parse data
		// data is in the format of "usernameNULLpasswordNULL"
		// username and password are separated by NULL character
		// NULL character is represented by byte 0x00
		parts := bytes.Split(data, []byte{0x00})
		if len(parts) != 2 {
			fmt.Println("Invalid data format")
			return
		}
		username := string(parts[0])
		password := string(parts[1])
		fmt.Printf("Username: %s, Password: %s\n", username, password)

	}
}

func main() {
	protocolVersion := "0001"
	messageType := "INIT"
	data := "user1\x00password1"
	length := fmt.Sprintf("%04d", 4+4+4+len(data))

	message := []byte(length + protocolVersion + messageType + data)

	ParseRequest(message)
}

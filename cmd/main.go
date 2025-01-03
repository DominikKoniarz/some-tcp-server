package main

import (
	"bytes"
	"fmt"
	"strconv"
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
	messageType := rawData[8:9]
	// rest of the bytes for data
	data := rawData[9:]

	if string(protocolVersion) != "0001" {
		fmt.Println("Invalid protocol version")
		return
	}

	_, err := strconv.Atoi(string(length))
	if err != nil {
		fmt.Println("Error converting length to integer:", err)
		return
	}

	stringMessageType := string(messageType)

	if stringMessageType == "I" {
		// split data by 0x00 to get username and password
		parts := bytes.Split(data, []byte{0x00})
		if len(parts) != 2 {
			fmt.Println("Invalid data format")
			return
		}
		username := string(parts[0])
		password := string(parts[1])

		fmt.Printf("Username: %s, Password: %s\n", username, password)
	} else {
		fmt.Println("Invalid message type")
	}
}

func main() {
	protocolVersion := "0001"
	messageType := "I"
	data := "user1\x00password1"
	length := fmt.Sprintf("%04d", 4+4+1+len(data))

	message := []byte(length + protocolVersion + messageType + data)

	ParseRequest(message)
}

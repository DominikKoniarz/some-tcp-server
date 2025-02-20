package request

import (
	"errors"
	"fmt"
)

const MAX_REQUEST_LENGTH = 1024

const AUTH_MESSAGE_TYPE = "A"

type Request struct {
	ProtocolVersion string
	MessageType     string // A - auth
	Data            string // change to byte?
}

func BuildRequest(messageType, data string) (Request, error) {
	protocolVersion := "0001"

	if len(messageType) != 1 {
		return Request{}, errors.New("invalid message type")
	}

	if (len(messageType) + len(data)) == 0 {
		return Request{}, errors.New("message type and data cannot be empty")
	}

	totalLength := len(protocolVersion) + len(messageType) + len(data)

	if totalLength > MAX_REQUEST_LENGTH {
		return Request{}, errors.New("data too long")
	}

	request := Request{
		ProtocolVersion: protocolVersion,
		MessageType:     messageType,
		Data:            data,
	}

	return request, nil
}

func (r Request) Prettify() string {
	return fmt.Sprintf("Request:\n  Protocol Version: %s\n  Message Type: %s\n  Data: %s", r.ProtocolVersion, r.MessageType, r.Data)
}

func (r Request) ToBytes() []byte {
	return []byte(r.ProtocolVersion + r.MessageType + r.Data)
}

func ParseRequest(rawData []byte) (Request, error) {
	if len(rawData) < 5 {
		return Request{}, errors.New("invalid data format")
	}

	protocolVersion := rawData[:4]
	messageType := rawData[4:5]
	data := rawData[5:]

	if string(protocolVersion) != "0001" {
		return Request{}, errors.New("invalid protocol version")
	}

	return Request{
		ProtocolVersion: string(protocolVersion),
		MessageType:     string(messageType),
		Data:            string(data),
	}, nil
}

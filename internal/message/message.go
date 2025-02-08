package message

import (
	"errors"
	"strings"
)

type AuthMessage struct {
	Username string
	Password string
}

// username:password
func (m AuthMessage) ToBytes() []byte {
	return []byte(m.Username + ":" + m.Password)
}

func ParseAuthMessage(data string) (AuthMessage, error) {
	// username:password, min length 3 (a:b)
	if len(data) < 3 {
		return AuthMessage{}, errors.New("invalid data format")
	}

	split := strings.Split(data, ":")
	if len(split) != 2 {
		return AuthMessage{}, errors.New("invalid data format")
	}

	return AuthMessage{
		Username: split[0],
		Password: split[1],
	}, nil
}

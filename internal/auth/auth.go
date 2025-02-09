package auth

import (
	"errors"

	"github.com/DominikKoniarz/some-tcp-server/internal/connection"
	"github.com/DominikKoniarz/some-tcp-server/internal/message"
	"github.com/DominikKoniarz/some-tcp-server/internal/request"
)

const USERNAME = "root"
const PASSWORD = "root"

func HandleAuth(c *connection.Connection, r *request.Request) error {
	if r.MessageType != request.AUTH_MESSAGE_TYPE {
		return errors.New("invalid message type")
	}

	authMessage, err := message.ParseAuthMessage(r.Data)
	if err != nil {
		return err
	}

	if authMessage.Username == USERNAME && authMessage.Password == PASSWORD {
		c.IsAuthenticated = true
		if _, err := (*c.C).Write([]byte("Authenticated")); err != nil {
			return err
		}
		return nil
	} else {
		if _, err := (*c.C).Write([]byte("Invalid credentials")); err != nil {
			return err
		}
		return errors.New("invalid credentials")
	}
}

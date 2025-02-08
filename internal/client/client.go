package client

import (
	"fmt"
	"net"

	"github.com/DominikKoniarz/some-tcp-server/internal/env"
)

type Client struct {
	c *net.Listener
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read data from connection
	for {
		buffer := make([]byte, 1024)

		_, err := conn.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Connection closed by client:", conn.RemoteAddr())
				break
			}

			fmt.Println("Error reading data:", err)
			break
		}
	}

}

func (c *Client) Connect() error {
	envs := env.LoadClientEnvs()

	address := fmt.Sprintf("%s:%s", envs.Host, envs.Port)

	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return err
	}

	c.handleConnection(conn)

	return nil
}

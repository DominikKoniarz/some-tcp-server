package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/DominikKoniarz/some-tcp-server/internal/env"
	"github.com/DominikKoniarz/some-tcp-server/internal/request"
)

type Client struct {
	C               *net.Listener
	IsAuthenticated bool
}

func NewClient() *Client {
	return &Client{IsAuthenticated: false}
}

// send credentials
func (c *Client) sendCredentials(conn net.Conn, env env.ClientEnvs) error {
	request, err := request.BuildRequest(request.AUTH_MESSAGE_TYPE, env.Username+":"+env.Password)
	if err != nil {
		fmt.Println("Error building request:", err)
		return err
	}

	_, err = conn.Write(request.ToBytes())
	if err != nil {
		fmt.Println("Error sending credentials:", err)
		return err
	}

	return nil
}

func (c *Client) handleConnection(conn net.Conn, env env.ClientEnvs) {
	defer func() {
		fmt.Println("Closing connection...")
		conn.Close()
	}()

	scanner := bufio.NewScanner(os.Stdin)

	var isAuthenticated bool

	buf := make([]byte, 1024)

	for {
		if !isAuthenticated {
			err := c.sendCredentials(conn, env)
			if err != nil {
				fmt.Println("Error sending credentials:", err)
				return
			}
		} else {
			fmt.Print("Enter text: ")
			scanner.Scan()
			text := scanner.Text()

			if text == "" {
				fmt.Println("No input received, please enter some text.")
				continue
			}

			if text == "exit" {
				fmt.Println("Exiting...")
				break
			}

			request, err := request.BuildRequest("1", text)
			if err != nil {
				fmt.Println("Error building request:", err)
				return
			}

			_, err = conn.Write(request.ToBytes())
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}
		}

		// usefull if we want to await particular response
		// wait for response for 5 seconds max
		// err := conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		// if err != nil {
		// 	fmt.Println("Error setting read deadline:", err)
		// 	return
		// }

		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading response:", err)

			if err.Error() == "EOF" {
				fmt.Println("Connection closed by server")
				break
			}
			return
		}

		fmt.Println("Received response:", string(buf[:n]))

		if string(buf[:n]) == "Authenticated" {
			isAuthenticated = true
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

	var conn *net.TCPConn
	for i := 0; i < 3; i++ {
		conn, err = net.DialTCP("tcp", nil, addr)
		if err == nil {
			break
		}
		fmt.Println("Error connecting to server, retrying...", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		fmt.Println("Failed to connect to server after 3 attempts:", err)
		return err
	}

	fmt.Println("Connected to server:", conn.RemoteAddr())

	c.handleConnection(conn, envs)

	return nil
}

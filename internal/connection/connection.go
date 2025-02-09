package connection

import "net"

type Connection struct {
	ID              string
	IsAuthenticated bool
	C               *net.Conn
}
